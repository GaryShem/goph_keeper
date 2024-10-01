package interceptors

import (
	"context"

	"google.golang.org/grpc"

	"goph_keeper/goph_server/internal/logging"
)

type MethodLogInterceptor struct {
}

func (a MethodLogInterceptor) Intercept(
	ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler,
) (interface{}, error) {
	logging.Log().Info(info.FullMethod)
	res, err := handler(ctx, req)
	if err != nil {
		logging.Log().Info(err.Error())
	} else {
		logging.Log().Info("no error")
	}
	return res, err
}
