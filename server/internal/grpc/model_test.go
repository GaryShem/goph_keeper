package grpc

import (
	"context"
	"fmt"
	"log"
	"net"
	"strconv"
	"testing"
	"time"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"goph_keeper/goph_server/internal/config"
	"goph_keeper/goph_server/internal/storage/postgresql"
	"goph_keeper/goph_server/internal/storage/repo"
	"goph_keeper/shared/proto"
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
		repository, err := postgresql.NewPostgreSQLRepo(serverConfig)
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

func SetupGrpcServer(repo repo.Repo) (context.Context, context.CancelFunc, *grpc.Server) {
	ctx, cancel := context.WithCancel(context.Background())
	server := RunGrpcServer(ctx, repo)

	return ctx, cancel, server
}

func TestGrpcSuite(t *testing.T) {
	suite.Run(t, new(GrpcSuite))
}

type GrpcSuite struct {
	suite.Suite
	repo         repo.Repo
	config       config.ServerConfig
	dbResource   *dockertest.Resource
	pool         *dockertest.Pool
	serverCtx    context.Context
	serverCancel context.CancelFunc
	server       *grpc.Server
	agent        proto.GophKeeperClient
}

func (s *GrpcSuite) SetupSuite() {
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

	ctx, cancel, server := SetupGrpcServer(repository)
	s.serverCtx = ctx
	s.serverCancel = cancel
	s.server = server

	go func() {
		listener, err := net.Listen("tcp", fmt.Sprintf(":%d", s.config.GrpcPort))
		s.Require().NoError(err)
		s.Require().NotNil(server)
		_ = server.Serve(listener)
	}()
	conn, err := grpc.NewClient(fmt.Sprintf(":%d", s.config.GrpcPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	s.Require().NoError(err)
	s.agent = proto.NewGophKeeperClient(conn)
}

func (s *GrpcSuite) TearDownSuite() {
	defer func() {
		if err := s.pool.Purge(s.dbResource); err != nil {
			log.Fatalf("Could not purge resource: %s", err)
		}
	}()
	s.serverCancel()
}

func (s *GrpcSuite) getContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 30*time.Second)
}
