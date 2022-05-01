package user

import (
	"context"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"google.golang.org/grpc"
)

func Interceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	//log.Println("unary server interceptor")
	return handler(ctx, req)
}

type TokenAuth struct {
	Token string
}

func (t TokenAuth) GetRequestMetadata(ctx context.Context, in ...string) (map[string]string, error) {
	//log.Println(in)
	return map[string]string{
		"token": t.Token,
	}, nil
}

func (t TokenAuth) RequireTransportSecurity() bool {
	return false
}

func a() {
	//metautils.()
	grpc_auth.UnaryServerInterceptor(func(ctx context.Context) (context.Context, error) {
		return ctx, nil
	})
}
