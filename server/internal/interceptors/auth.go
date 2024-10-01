package interceptors

import (
	"context"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"goph_keeper/goph_server/internal/storage/repo"
)

const REGISTER_USER_METHOD = "RegisterUser"
const PING_METHOD = "Ping"

var ErrAccessDenied = status.Error(codes.PermissionDenied, "access denied")

type LoginCreds struct {
	Username string
	Password string
}

func NewAuthInterceptor(repo repo.Repo) *AuthInterceptor {
	return &AuthInterceptor{
		repo: repo,
	}
}

type AuthInterceptor struct {
	repo repo.Repo
}

func (a *AuthInterceptor) authorize(ctx context.Context) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ErrAccessDenied
	}
	if len(md["username"]) == 0 || len(md["password"]) == 0 {
		return ErrAccessDenied
	}
	return a.repo.LoginUser(md["username"][0], md["password"][0])
}

func (a *AuthInterceptor) Intercept(
	ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler,
) (interface{}, error) {
	method := info.FullMethod
	if strings.Contains(method, REGISTER_USER_METHOD) || strings.Contains(method, PING_METHOD) {
		return handler(ctx, req)
	}
	if err := a.authorize(ctx); err != nil {
		return ctx, status.Error(codes.Internal, err.Error())
	}
	return handler(ctx, req)
}
