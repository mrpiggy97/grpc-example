package main

import (
	"fmt"
	"net"
	"os"

	"github.com/mrpiggy97/grpcExample/interceptors"
	"github.com/mrpiggy97/sharedProtofiles/calculation"
	"github.com/mrpiggy97/sharedProtofiles/formatting"
	"github.com/mrpiggy97/sharedProtofiles/num"
	"github.com/mrpiggy97/sharedProtofiles/randomNumber"
	"github.com/mrpiggy97/sharedProtofiles/user"

	"google.golang.org/grpc"
)

func main() {
	var port string = os.Getenv("PORT")
	if len(port) == 0 {
		panic("port cannot be blank")
	}
	var address string = fmt.Sprintf("0.0.0.0:%v", port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		panic(err)
	}
	var userServer *user.Server = new(user.Server)
	var grpcServer *grpc.Server = grpc.NewServer(
		interceptors.WithMessageInterceptor(),
	)
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
