package main

import (
	"context"
	"log"
	"net"
	"net/http"

	"github.com/changpro/disk-service/gateway/interfaces"
	gatewaypb "github.com/changpro/disk-service/resource/stub"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Create a listener on TCP port
	lis, err := net.Listen("tcp", ":8001")
	if err != nil {
		log.Fatalln("Failed to listen:", err)
	}

	// Create a gRPC server object
	s := grpc.NewServer()
	// Attach the Greeter service to the server
	gatewaypb.RegisterGatewayServiceServer(s, interfaces.NewServer())
	// Serve gRPC Server
	log.Println("Serving gRPC on localhost:8001")
	go func() {
		log.Fatalln(s.Serve(lis))
	}()

	// Create a client connection to the gRPC server we just started
	// This is where the gRPC-Gateway proxies the requests
	conn, err := grpc.DialContext(
		context.Background(),
		"localhost:8001",
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalln("Failed to dial server:", err)
	}

	gwmux := runtime.NewServeMux()
	// Register RESTful api gateway
	err = gatewaypb.RegisterGatewayServiceHandler(context.Background(), gwmux, conn)
	if err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}

	gwServer := &http.Server{
		Addr:    ":8002",
		Handler: gwmux,
	}

	log.Println("Serving gRPC-Gateway on http://localhost:8002")
	log.Fatalln(gwServer.ListenAndServe())
}
