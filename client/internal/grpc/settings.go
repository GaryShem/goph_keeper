package grpc

import "strconv"

func (g *GrpcWrapper) GrpcSetHost(host string) {
	g.settings.Host = host
}
func (g *GrpcWrapper) GrpcSetPort(port string) error {
	portInt, err := strconv.Atoi(port)
	if err != nil {
		return err
	}
	g.settings.Port = portInt
	return nil
}
func (g *GrpcWrapper) GrpcSetName(name string) {
	g.settings.Username = name
}
func (g *GrpcWrapper) GrpcSetPassword(password string) {
	g.settings.Password = password
}

type ClientSettings struct {
	Username string
	Password string
	Host     string
	Port     int
}
