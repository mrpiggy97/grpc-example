package main

import (
	"fmt"
	"net"

	"github.com/mrpiggy97/sharedProtofiles/calculation"
	"github.com/mrpiggy97/sharedProtofiles/formatting"
	"github.com/mrpiggy97/sharedProtofiles/num"
	"github.com/mrpiggy97/sharedProtofiles/randomNumber"
	"github.com/mrpiggy97/sharedProtofiles/user"

	"google.golang.org/grpc"
)

func main() {
	listener, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		panic(err)
	}
	var userServer *user.Server = new(user.Server)
	var grpcServer *grpc.Server = grpc.NewServer()
	user.RegisterUserServiceServer(grpcServer, userServer)

	var numServer *num.Server = new(num.Server)
	num.RegisterNumServiceServer(grpcServer, numServer)

	var formattingServer *formatting.Server = new(formatting.Server)
	formatting.RegisterFormattingServiceServer(grpcServer, formattingServer)

	var randomNumberServer *randomNumber.Server = new(randomNumber.Server)
	randomNumber.RegisterRandomServiceServer(grpcServer, randomNumberServer)

	var calculationServer *calculation.Server = new(calculation.Server)
	calculation.RegisterCalculationServiceServer(grpcServer, calculationServer)

	fmt.Println("server is listening at port 50051")
	if listeningErr := grpcServer.Serve(listener); listeningErr != nil {
		panic(listeningErr)
	}
}
