package interceptors

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func LoggerInterceptor(
	cxt context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	metaInfo, found := metadata.FromIncomingContext(cxt)
	if !found {
		return nil, status.Errorf(codes.Unimplemented, "no logs were found")
	}
	os, osFound := metaInfo["os"]
	if !osFound {
		return nil, status.Errorf(codes.NotFound, "os not found")
	}
	if os[0] == "" {
		return nil, status.Errorf(codes.InvalidArgument, "os cannot be blank")
	}
	zone, zoneFound := metaInfo["zone"]
	if !zoneFound {
		return nil, status.Errorf(codes.NotFound, "zone not found")
	}

	if zone[0] == "" {
		return nil, status.Errorf(codes.InvalidArgument, "zone cannot be blank")
	}

	object, err := handler(cxt, req)
	return object, err
}

func WithLoggerInterceptor() grpc.ServerOption {
	return grpc.UnaryInterceptor(LoggerInterceptor)
}
