package postgresql

import (
	"fmt"
	"log"
	"strconv"
	"testing"
	"time"

	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/suite"

	"goph_keeper/goph_server/internal/config"
	"goph_keeper/goph_server/internal/storage/repo"

	"github.com/ory/dockertest/v3"
)

func SetupPostgresContainer(serverConfig config.ServerConfig) (*dockertest.Pool, *dockertest.Resource, repo.Repo) {
	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not construct pool: %s", err)
	}
	pool = pool

	// uses pool to try to connect to Docker
	err = pool.Client.Ping()
	if err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
	}

	// pulls an image, creates a container based on it and runs it
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "latest",
		Env: []string{
			fmt.Sprintf("POSTGRES_PASSWORD=%s", serverConfig.DbPass),
			fmt.Sprintf("POSTGRES_USER=%s", serverConfig.DbUser),
			fmt.Sprintf("POSTGRES_DB=%s", serverConfig.DbName),
			"listen_addresses = '*'",
		},
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}
	_ = resource.Expire(300)

	hostAndPort := resource.GetHostPort("5432/tcp")
	port, err := strconv.Atoi(resource.GetPort("5432/tcp"))
	if err != nil {
		log.Fatalf("Could not convert port to int: %s", err)
	}
	serverConfig.DbPort = port
	databaseUrl := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		serverConfig.DbUser, serverConfig.DbPass, hostAndPort, serverConfig.DbName)

	log.Println("Connecting to database on url: ", databaseUrl)

	var dockerTimeout uint = 60
	// Tell docker to hard kill the container after a timeout
	if err = resource.Expire(dockerTimeout); err != nil {
		log.Fatalf("Could not expire resource: %s", err)
	}

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	pool.MaxWait = time.Duration(dockerTimeout) * time.Second
	var returnRepo repo.Repo
	if err = pool.Retry(func() error {
		repository, err := NewPostgreSQLRepo(serverConfig)
		if err != nil {
			return err
		}
		returnRepo = repository
		return nil
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}
	_ = returnRepo.RegisterUser("test", "test")
	return pool, resource, returnRepo
}

func TestPostgreSQLSuite(t *testing.T) {
	suite.Run(t, new(PostgreSQLSuite))
}

type PostgreSQLSuite struct {
	suite.Suite
	repo       repo.Repo
	config     config.ServerConfig
	dbResource *dockertest.Resource
	pool       *dockertest.Pool
}

func (s *PostgreSQLSuite) SetupSuite() {
	s.config = config.ServerConfig{
		DbHost:   "localhost",
		DbPort:   5432,
		DbUser:   "test",
		DbPass:   "test",
		DbName:   "goph",
		GrpcPort: 8080,
	}
	pool, resource, repository := SetupPostgresContainer(s.config)
	s.pool = pool
	s.dbResource = resource
	s.repo = repository
}

func (s *PostgreSQLSuite) TearDownSuite() {
	defer func() {
		if err := s.pool.Purge(s.dbResource); err != nil {
			log.Fatalf("Could not purge resource: %s", err)
		}
	}()
}

func (s *PostgreSQLSuite) TestPing() {
	err := s.repo.Ping()
	s.Require().NoError(err)
}
