package formatting_test

import (
	"context"
	"fmt"
	"net"
	"sync"
	"testing"

	"github.com/mrpiggy97/sharedProtofiles/formatting"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

var listener *bufconn.Listener = bufconn.Listen(1024 * 1024)

func bufDialer(cxt context.Context, str string) (net.Conn, error) {
	return listener.Dial()
}

func runServer(stopServer *sync.WaitGroup) {
	var grpcServer *grpc.Server = grpc.NewServer()
	var formattingServer *formatting.Server = new(formatting.Server)
	formatting.RegisterFormattingServiceServer(grpcServer, formattingServer)
	grpcServer.Serve(listener)
	stopServer.Wait()
	defer listener.Close()
}

func runClient(sendClient chan<- formatting.FormattingServiceClient) {
	conn, connError := grpc.DialContext(
		context.Background(),
		"bufnet",
		grpc.WithInsecure(),
		grpc.WithContextDialer(bufDialer),
	)
	if connError != nil {
		panic("failed to establish testing connection between client and server")
	}

	var client formatting.FormattingServiceClient = formatting.NewFormattingServiceClient(conn)
	sendClient <- client
}

func TestToLowerCase(testCase *testing.T) {
	//run testing servers
	var stopServer *sync.WaitGroup = new(sync.WaitGroup)
	stopServer.Add(1)
	var getClient chan formatting.FormattingServiceClient = make(chan formatting.FormattingServiceClient, 1)
	go runServer(stopServer)
	go runClient(getClient)

	//get client
	var client formatting.FormattingServiceClient = <-getClient

	//make request
	var request *formatting.FormattingRequest = &formatting.FormattingRequest{
		StringToConvert: "THS-IS-THE-STRING",
	}

	res, resError := client.ToLowerCase(
		context.Background(),
		request,
	)

	//run tests
	if resError != nil {
		testCase.Error(resError)
	}

	var expectedResponse string = "thsisthestring"
	if res.GetConvertedString() != expectedResponse {
		message := fmt.Sprintf("expected res.ConvertedString to be %v, instead got %v", expectedResponse, res.ConvertedString)
		testCase.Error(message)
	}

	defer stopServer.Done()
}

func TestToCamelCase(testCase *testing.T) {
	//run testing servers
	var stopServer *sync.WaitGroup = new(sync.WaitGroup)
	stopServer.Add(1)
	var getClient chan formatting.FormattingServiceClient = make(chan formatting.FormattingServiceClient, 1)
	go runServer(stopServer)
	go runClient(getClient)

	//get client
	var client formatting.FormattingServiceClient = <-getClient

	//make request
	var request *formatting.FormattingRequest = &formatting.FormattingRequest{
		StringToConvert: "my-name-is",
	}
	res, resError := client.ToCamelCase(
		context.Background(),
		request,
	)

	//run tests
	if resError != nil {
		testCase.Error(resError)
	}

	var expectedResponse string = "MyNameIs"
	if res.GetConvertedString() != expectedResponse {
		message := fmt.Sprintf("expected res.ConvertedString to be %v, instead got %v",
			expectedResponse, res.GetConvertedString())
		testCase.Error(message)
	}
	defer stopServer.Done()
}

func TestToUpperCase(testCase *testing.T) {
	//run testing servers
	var stopServer *sync.WaitGroup = new(sync.WaitGroup)
	stopServer.Add(1)
	var getClient chan formatting.FormattingServiceClient = make(chan formatting.FormattingServiceClient, 1)

	go runServer(stopServer)
	go runClient(getClient)

	//get client
	var client formatting.FormattingServiceClient = <-getClient

	//make request
	var request *formatting.FormattingRequest = &formatting.FormattingRequest{
		StringToConvert: "convert-this-coon-yo",
	}
	res, resErr := client.ToUpperCase(
		context.Background(),
		request,
	)

	//make tests
	if resErr != nil {
		testCase.Error(resErr)
	}

	var expectedResponse string = "CONVERTTHISCOONYO"

	if res.GetConvertedString() != expectedResponse {
		message := fmt.Sprintf("expected res.ConvertedString to be %v, instead got %v",
			expectedResponse, res.ConvertedString)
		testCase.Error(message)
	}
}
