package interceptors

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func GeneralInterceptor(
	cxt context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	metaInfo, metaInfoFound := metadata.FromIncomingContext(cxt)
	if metaInfoFound {
		fmt.Println(metaInfo)
	}
	authority, authorityFound := metaInfo[":authority"]
	if !authorityFound {
		return nil, status.Errorf(codes.PermissionDenied, "request has not given :authority header")
	}
	fmt.Println(authority[0])
	return handler(cxt, req)
}

func WithGeneralInterceptor() grpc.ServerOption {
	return grpc.UnaryInterceptor(GeneralInterceptor)
}
