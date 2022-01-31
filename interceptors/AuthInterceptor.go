package interceptors

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func AuthInterceptor(
	cxt context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	metaInfo, metaDataFound := metadata.FromIncomingContext(cxt)
	if !metaDataFound {
		return nil, status.Errorf(codes.InvalidArgument, "argument not found")
	}
	password, passwordFound := metaInfo["password"]
	if !passwordFound {
		return nil, status.Errorf(codes.Unauthenticated, "password not found")
	}
	if password[0] != "go" {
		return nil, status.Errorf(codes.InvalidArgument, "password not valid")
	}
	cxt = metadata.AppendToOutgoingContext(cxt, "is_authenticated", "true")
	object, err := handler(cxt, req)
	return object, err
}

func WithAuthInterceptor() grpc.ServerOption {
	return grpc.UnaryInterceptor(AuthInterceptor)
}
