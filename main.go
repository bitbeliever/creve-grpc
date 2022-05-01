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

	grpcClientTest()

	httpClient()
}

func grpcClientTest() {
	conn, err := grpc.Dial("0.0.0.0:9090",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		panic(err)
	}

	cli := userservicepb.NewUserServiceClient(conn)
	resp, err := cli.Test(context.Background(), &userservicepb.TestRequest{
		Name:  "plopski",
		Age:   "20",
		Query: "query",
		Token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6MSwiQ3JlYXRlZEF0IjoiMjAyMi0wNC0yM1QyMjozNToyMi4yMyswODowMCIsIlVwZGF0ZWRBdCI6IjIwMjItMDQtMjNUMjI6MzU6MjIuMjMrMDg6MDAiLCJEZWxldGVkQXQiOm51bGwsIlVzZXJuYW1lIjoic2hveCIsIlBhc3N3b3JkIjoiMTIzNCIsIkF2YXRhciI6Imh0dHBzOi8vZ3cuYWxpcGF5b2JqZWN0cy5jb20vem9zL2FudGZpbmNkbi9YQW9zWHVOWnlGL0JpYXpmYW54bWFtTlJveHhWeGthLnBuZyIsIkVtYWlsIjoiIiwiU2FsdCI6IiIsIlBob25lIjoiIiwiU3RhdHVzIjowfQ.V8BSInvjIdAJVdsm1zJOFqpPcA75rkHhJg6wwM_BvC8",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("client test:", resp)
}

func httpClient() {
	conn, err := grpc.DialContext(
		context.Background(),
		"0.0.0.0:9090",
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
		Addr:    "0.0.0.0:9091",
		Handler: gwmux,
	}
	log.Println("Serving gRPC-Gateway on http://0.0.0.0:9091")
	log.Fatalln(gwServer.ListenAndServe())
}
