package main

import (
	"fmt"
	"net"

	"github.com/mrpiggy97/grpc-example/user"
	"google.golang.org/grpc"
)

func main() {
	listener, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		panic(err)
	}
	var grpcServer *grpc.Server = grpc.NewServer()
	var userServerInstance *user.UserServer = new(user.UserServer)
	user.RegisterUserServiceServer(grpcServer, userServerInstance)
	fmt.Println("server is listening at port 50051")
	if listeningErr := grpcServer.Serve(listener); listeningErr != nil {
		panic(listeningErr)
	}
}
