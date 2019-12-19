package middleware

import (
	"context"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/ezio1119/fishapp-post/conf"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	create string = "/post_grpc.PostService/Create"
	update string = "/post_grpc.PostService/Update"
	delete string = "/post_grpc.PostService/Delete"
)

func (*middleware) AuthInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		var err error
		method := info.FullMethod

		if method == create || method == update || method == delete {
			ctx, err = authFunc(ctx)
		}
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, err.Error())
		}
		return handler(ctx, req)
	}
}

func authFunc(ctx context.Context) (context.Context, error) {
	t, err := grpc_auth.AuthFromMD(ctx, "bearer")
	if err != nil {
		return nil, err
	}
	userID, err := parseToken(t)
	if err != nil {
		return nil, err
	}
	newCtx := context.WithValue(ctx, "userID", userID)
	return newCtx, nil
}

func parseToken(t string) (int64, error) {
	jwtkey := []byte(conf.C.Auth.Jwtkey)
	var claims jwt.StandardClaims
	_, err := jwt.ParseWithClaims(t, &claims, func(token *jwt.Token) (interface{}, error) {
		return jwtkey, nil
	})
	if err != nil {
		return 0, err
	}
	userID, err := strconv.ParseInt(claims.Subject, 10, 64)
	if err != nil {
		return 0, err
	}
	return userID, nil
}
