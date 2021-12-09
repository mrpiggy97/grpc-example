package user

import (
	"context"
	"fmt"
)

type UserServer struct {
	UnimplementedUserServiceServer
}

func (userServerInstance *UserServer) GetUser(cxt context.Context, req *UserRequest) (*User, error) {
	fmt.Println("server recieved:", req.String())
	var newUser *User = &User{UserId: "John", Email: "example@email.com"}
	return newUser, nil
}
