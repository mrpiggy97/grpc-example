package num_test

import (
	"context"
	"fmt"
	"io"
	"net"
	"sync"
	"testing"

	"github.com/mrpiggy97/sharedProtofiles/num"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

var listener *bufconn.Listener = bufconn.Listen(1024 * 1024)

func bufDialer(cxt context.Context, str string) (net.Conn, error) {
	return listener.Dial()
}

func runServer(stopServer *sync.WaitGroup) {
	var grpcServer *grpc.Server = grpc.NewServer()
	var numServer *num.Server = new(num.Server)
	num.RegisterNumServiceServer(grpcServer, numServer)
	grpcServer.Serve(listener)
	stopServer.Wait()
	defer listener.Close()
}

func runClient(sendClient chan<- num.NumServiceClient) {
	conn, connError := grpc.DialContext(
		context.Background(),
		"bufnet",
		grpc.WithInsecure(),
		grpc.WithContextDialer(bufDialer),
	)
	if connError != nil {
		panic("failed to establish testing connection between client and server")
	}

	var client num.NumServiceClient = num.NewNumServiceClient(conn)
	sendClient <- client
}

func TestRnd(testCase *testing.T) {
	//run test servers
	var getClient chan num.NumServiceClient = make(chan num.NumServiceClient, 1)
	var stopServer *sync.WaitGroup = new(sync.WaitGroup)
	stopServer.Add(1)
	go runServer(stopServer)
	go runClient(getClient)

	//get client
	var client num.NumServiceClient = <-getClient

	//make request
	var request *num.NumRequest = &num.NumRequest{
		From:   0,
		To:     50,
		Number: 67,
	}

	stream, streamError := client.Rnd(
		context.Background(),
		request,
	)

	stream.CloseSend()

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

func TestSum(testCase *testing.T) {
	//run testing servers
	var getClient chan num.NumServiceClient = make(chan num.NumServiceClient, 1)
	var stopServer *sync.WaitGroup = new(sync.WaitGroup)
	stopServer.Add(1)
	go runServer(stopServer)
	go runClient(getClient)

	//get client
	var client num.NumServiceClient = <-getClient

	stream, streamError := client.Sum(
		context.Background(),
	)

	if streamError != nil {
		testCase.Error(streamError)
	}

	//send stream of requests
	for i := 0; i < 20; i++ {
		var request *num.SumRequest = &num.SumRequest{
			Number: int64(i),
		}
		var sendingErr error = stream.Send(request)
		if sendingErr != nil {
			testCase.Error(sendingErr)
		}
	}
	res, resError := stream.CloseAndRecv()
	if resError != nil {
		testCase.Error(resError)
	}
	fmt.Println(res.String())
	defer stopServer.Done()
}
