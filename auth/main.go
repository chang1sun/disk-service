package main

import (
	"log"
	"net"

	authpb "github.com/changpro/disk-service/auth/deps"
	"github.com/changpro/disk-service/auth/interfaces"
	"google.golang.org/grpc"
)

func main() {
	// Create a listener on TCP port
	lis, err := net.Listen("tcp", ":8003")
	if err != nil {
		log.Fatalln("Failed to listen:", err)
	}

	// Create a gRPC server object
	s := grpc.NewServer()
	// Attach the Greeter service to the server
	authpb.RegisterAuthServiceServer(s, interfaces.NewServer())
	// Serve gRPC Server
	log.Println("Serving gRPC on localhost:8003")
	log.Fatalln(s.Serve(lis))
}
