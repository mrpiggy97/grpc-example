package user_test

import (
	"context"
	"fmt"
	"io"
	"net"
	"sync"
	"testing"

	"github.com/mrpiggy97/sharedProtofiles/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/test/bufconn"
)

type clientWrapper struct {
	conn   *grpc.ClientConn
	client user.UserServiceClient
}

var listener *bufconn.Listener = bufconn.Listen(1024 * 1024)

func buffDialer(cxt context.Context, str string) (net.Conn, error) {
	return listener.Dial()
}

func runServer(stopServer *sync.WaitGroup) {
	var grpcServer *grpc.Server = grpc.NewServer()
	var userServer *user.Server = new(user.Server)
	user.RegisterUserServiceServer(grpcServer, userServer)
	grpcServer.Serve(listener)
	stopServer.Wait()
	defer listener.Close()
}

func runClient(sendWrapper chan<- clientWrapper) {
	connection, connErr := grpc.DialContext(
		context.Background(),
		"buffnet",
		grpc.WithInsecure(),
		grpc.WithContextDialer(buffDialer),
	)
	if connErr != nil {
		panic("error trying to connect server and client")
	}
	var userClient user.UserServiceClient = user.NewUserServiceClient(connection)
	var wrapper clientWrapper = clientWrapper{
		conn:   connection,
		client: userClient,
	}
	sendWrapper <- wrapper
}

func TestGetUser(testCase *testing.T) {
	//run testing servers
	var waiter *sync.WaitGroup = new(sync.WaitGroup)
	waiter.Add(1)
	var getClient chan clientWrapper = make(chan clientWrapper, 1)
	go runServer(waiter)
	go runClient(getClient)

	//get client
	var client clientWrapper = <-getClient

	//make request
	var request *user.UserRequest = &user.UserRequest{
		UserId: "1223123wasdasds",
	}
	var cxt context.Context = context.Background()
	cxt = metadata.AppendToOutgoingContext(cxt, "password", "go")
	res, resError := client.client.GetUser(
		cxt,
		request,
	)
	if resError != nil {
		testCase.Error(resError)
	}
	fmt.Println(res.String())
	defer waiter.Done()
}

func TestRegisterUsers(testCase *testing.T) {
	//run servers
	var stopServer *sync.WaitGroup = new(sync.WaitGroup)
	stopServer.Add(1)
	var getClientWrapper chan clientWrapper = make(chan clientWrapper, 1)
	go runServer(stopServer)
	go runClient(getClientWrapper)

	//get client
	var client clientWrapper = <-getClientWrapper
	stream, streamErr := client.client.RegisterUsers(
		context.Background(),
	)

	if streamErr != nil {
		testCase.Error("error getting stream to register users")
	}

	//make request
	var baseUsername string = "cochinito"
	for i := 0; i < 10; i++ {
		var request *user.RegisterUserRequest = &user.RegisterUserRequest{
			Username: fmt.Sprintf("%v %v", baseUsername, i),
		}
		var sendingError error = stream.Send(request)
		if sendingError != nil {
			testCase.Error(sendingError)
		}
	}
	stream.CloseSend()
	for {
		response, resError := stream.Recv()
		if resError != nil && resError != io.EOF {
			testCase.Error("resError should only be io.EOF")
		}
		if resError == io.EOF {
			fmt.Println("finished consuming stream")
			break
		}
		fmt.Println(response.String())
	}
	defer stopServer.Done()
}
