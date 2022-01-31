package randomNumber_test

import (
	"context"
	"fmt"
	"io"
	"net"
	"sync"
	"testing"

	"github.com/mrpiggy97/sharedProtofiles/randomNumber"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

var listener *bufconn.Listener = bufconn.Listen(1024 * 1024)

func BuffDialer(cxt context.Context, str string) (net.Conn, error) {
	return listener.Dial()
}

func runServer(stopServer *sync.WaitGroup) {
	var grpcServer *grpc.Server = grpc.NewServer()
	var randomNumberServer *randomNumber.Server = new(randomNumber.Server)
	randomNumber.RegisterRandomServiceServer(grpcServer, randomNumberServer)
	grpcServer.Serve(listener)
	stopServer.Wait()
	defer listener.Close()
}

func runClient(sendClient chan<- randomNumber.RandomServiceClient) {
	connection, connErr := grpc.DialContext(
		context.Background(),
		"bufnet",
		grpc.WithInsecure(),
		grpc.WithContextDialer(BuffDialer),
	)
	if connErr != nil {
		panic("failed to create connection between testing server")
	}
	var client randomNumber.RandomServiceClient = randomNumber.NewRandomServiceClient(connection)
	sendClient <- client
}

func TestAddRandomNumber(testCase *testing.T) {
	//run testing servers
	var stopServer *sync.WaitGroup = new(sync.WaitGroup)
	stopServer.Add(1)
	var getClient chan randomNumber.RandomServiceClient = make(chan randomNumber.RandomServiceClient, 1)
	go runServer(stopServer)
	go runClient(getClient)

	//get client
	var client randomNumber.RandomServiceClient = <-getClient

	//make request
	var request *randomNumber.RandomNumberRequest = &randomNumber.RandomNumberRequest{
		Number: 545482,
	}
	stream, streamError := client.AddRandomNumber(
		context.Background(),
		request,
	)
	if streamError != nil {
		testCase.Error(streamError)
	}
	stream.CloseSend()

	for {
		res, resError := stream.Recv()
		if resError != nil && resError != io.EOF {
			testCase.Error("resError should only be io.EOF,instead got ", resError)
		}
		if resError == io.EOF {
			fmt.Println("finished consuming stream")
			break
		}
		fmt.Println(res.String())
	}
	defer stopServer.Done()
}

func TestSubstractRandomNumber(testCase *testing.T) {
	//run testing servers
	var stopServer *sync.WaitGroup = new(sync.WaitGroup)
	stopServer.Add(1)
	var getClient chan randomNumber.RandomServiceClient = make(chan randomNumber.RandomServiceClient, 1)
	go runServer(stopServer)
	go runClient(getClient)

	//get client
	var client randomNumber.RandomServiceClient = <-getClient

	//make request
	var request *randomNumber.RandomNumberRequest = &randomNumber.RandomNumberRequest{
		Number: 123213213,
	}

	stream, streamError := client.SubstractRandomNumber(
		context.Background(),
		request,
	)
	if streamError != nil {
		testCase.Error(streamError)
	}

	//consume stream
	for {
		res, resError := stream.Recv()
		if resError != nil && resError != io.EOF {
			testCase.Error("resError should only be io.EOF, instead got ", resError)
		}
		if resError == io.EOF {
			fmt.Println("finished consuming stream")
			break
		}
		fmt.Println(res.String())
	}
	defer stopServer.Done()
}
