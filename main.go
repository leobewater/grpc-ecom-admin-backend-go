package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"strconv"

	"github.com/leobewater/udemy-orders-go-admin/database"
	"github.com/leobewater/udemy-orders-go-admin/models"
	"github.com/leobewater/udemy-orders-go-admin/pb/userpb"
	"github.com/leobewater/udemy-orders-go-admin/util"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type server struct{}

func (*server) User(ctx context.Context, req *userpb.UserLoginRequest) (*userpb.UserLoginResponse, error) {
	email := req.GetUser().GetEmail()
	password := req.GetUser().GetPassword()

	// find user
	var modelUser models.User
	// if user is found, map db data to &user
	database.DB.Where("email = ?", email).First(&modelUser)
	if modelUser.Id == 0 {
		res := &userpb.UserLoginResponse{
			Result: "User not found",
		}
		return res, nil
	}

	// compare password
	if err := modelUser.ComparePassword(password); err != nil {
		res := &userpb.UserLoginResponse{
			Result: "incorrect password",
		}
		return res, nil
	}

	// generate jwt and uses user.ID as the issuer
	token, err := util.GenerateJwt(strconv.Itoa(int(modelUser.Id)))
	if err != nil {
		res := &userpb.UserLoginResponse{
			Result: "Server error not able to generate jwt",
		}
		return res, nil
	}

	res := &userpb.UserLoginResponse{
		Result: token,
	}
	return res, nil
}

func main() {
	database.Connect()

	fmt.Println("Admin GRPC Server is running")

	// grpc port 50051
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// tls to enable or disable SSL
	tls := false
	opts := []grpc.ServerOption{}
	if tls {
		// SSL setup
		certFile := "ssl/server.crt"
		keyFile := "ssl/server.pem"
		creds, sslErr := credentials.NewServerTLSFromFile(certFile, keyFile)
		if sslErr != nil {
			log.Fatalf("Failed loading ssl certificates: %v", sslErr)
			return
		}
		opts = append(opts, grpc.Creds(creds))
	}

	// create grpc server
	s := grpc.NewServer(opts...)
	userpb.RegisterUserLoginServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
