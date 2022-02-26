package interceptors

import (
	"context"
	"fmt"

	"github.com/TwiN/go-color"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func MessageInterceptor(
	cxt context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	md, mdFound := metadata.FromIncomingContext(cxt)
	if !mdFound {
		return nil, status.Errorf(codes.NotFound, "metadata not found")
	}

	message, messageFound := md["message"]
	if !messageFound {
		return nil, status.Errorf(codes.NotFound, "message not found")
	}
	if message[0] == "" {
		return nil, status.Errorf(codes.InvalidArgument, "message cannot be empty or have length 0")
	}
	fmt.Println(color.Colorize(color.Red, message[0]))
	return handler(cxt, req)
}

func WithMessageInterceptor() grpc.ServerOption {
	return grpc.UnaryInterceptor(MessageInterceptor)
}
