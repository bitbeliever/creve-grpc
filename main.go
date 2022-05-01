package main

import (
	"context"
	"fmt"
	_ "github.com/bitbeliever/creve-grpc/configs"
	userservicepb "github.com/bitbeliever/creve-grpc/proto/user/pb"
	"github.com/bitbeliever/creve-grpc/service/user"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
	"net/http"
)

func main() {
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(user.Interceptor),
	)
	//userservicepb.RegisterUserServiceServer(grpcServer, &server{})
	userservicepb.RegisterUserServiceServer(grpcServer, &user.Service{})

	lis, err := net.Listen("tcp", ":9090")
	if err != nil {
		panic(err)
	}

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			panic(err)
		}
	}()

	//grpcClientTest()

	httpClient()
}

func grpcClientTest() {
	conn, err := grpc.Dial("localhost:9090",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		panic(err)
	}

	cli := userservicepb.NewUserServiceClient(conn)
	resp, err := cli.NewUser(context.Background(), &userservicepb.NewUserRequest{
		Username: "videogamedunkey",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(resp)
}

func httpClient() {
	conn, err := grpc.DialContext(
		context.Background(),
		"localhost:9090",
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithPerRPCCredentials(user.TokenAuth{Token: "token is sth"}),
	)
	if err != nil {
		log.Fatalln("Failed to dial server:", err)
	}
	gwmux := runtime.NewServeMux()
	err = userservicepb.RegisterUserServiceHandler(context.Background(), gwmux, conn)
	if err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}

	gwServer := &http.Server{
		Addr:    "localhost:9091",
		Handler: gwmux,
	}
	log.Println("Serving gRPC-Gateway on http://localhost:9091")
	log.Fatalln(gwServer.ListenAndServe())
}
