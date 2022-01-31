package calculation_test

import (
	"context"
	"fmt"
	"io"
	"net"
	"sync"
	"testing"

	"github.com/mrpiggy97/sharedProtofiles/calculation"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

var listener *bufconn.Listener = bufconn.Listen(1024 * 1024)

func bufDialer(cxt context.Context, str string) (net.Conn, error) {
	return listener.Dial()
}

func runServer(stopServer *sync.WaitGroup) {
	var grpcServer *grpc.Server = grpc.NewServer()
	var calculationServer *calculation.Server = new(calculation.Server)
	calculation.RegisterCalculationServiceServer(grpcServer, calculationServer)
	grpcServer.Serve(listener)
	stopServer.Wait()
	defer listener.Close()
}

func runClient(sendClient chan<- calculation.CalculationServiceClient) {
	conn, connErr := grpc.DialContext(
		context.Background(),
		"bufnet",
		grpc.WithInsecure(),
		grpc.WithContextDialer(bufDialer),
	)

	if connErr != nil {
		panic("error setting up connection between test servers")
	}

	var client calculation.CalculationServiceClient = calculation.NewCalculationServiceClient(conn)
	sendClient <- client
}

func TestSumStream(testCase *testing.T) {
	//run test servers
	var stopServer *sync.WaitGroup = new(sync.WaitGroup)
	stopServer.Add(1)
	var getClient chan calculation.CalculationServiceClient = make(chan calculation.CalculationServiceClient, 1)
	go runServer(stopServer)
	go runClient(getClient)

	//get client
	var client calculation.CalculationServiceClient = <-getClient

	//get stream
	stream, streamError := client.SumStream(context.Background())
	if streamError != nil {
		testCase.Error(streamError)
	}

	//make requests
	for i := 0; i < 20; i++ {
		var request *calculation.SumStreamRequest = &calculation.SumStreamRequest{
			A: int32(i),
			B: int32(i),
		}
		stream.Send(request)
	}
	stream.CloseSend()

	for {
		res, resError := stream.Recv()
		if resError != nil && resError != io.EOF {
			testCase.Error("resError shouold only be io.EOF, intead got ", resError)
		}
		if resError == io.EOF {
			fmt.Println("finished consuming stream")
			break
		}
		fmt.Println(res.String())
	}
	defer stopServer.Done()
}
