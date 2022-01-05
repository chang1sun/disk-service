package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	diskpb "github.com/changpro/disk-service/resource/stub"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type server struct {
	diskpb.UnimplementedFileServiceServer
}

func NewServer() *server {
	return &server{}
}

func (s *server) UploadFile(ctx context.Context, req *diskpb.UploadFileReq) (*diskpb.CommmonRsp, error) {
	if req.User == "" {
		return &diskpb.CommmonRsp{Status: "fail"}, fmt.Errorf("param invalid")
	}
	return &diskpb.CommmonRsp{Status: "success"}, nil
}

func main() {
	// Create a listener on TCP port
	lis, err := net.Listen("tcp", ":8001")
	if err != nil {
		log.Fatalln("Failed to listen:", err)
	}

	// Create a gRPC server object
	s := grpc.NewServer()
	// Attach the Greeter service to the server
	diskpb.RegisterFileServiceServer(s, &server{})
	// Serve gRPC Server
	log.Println("Serving gRPC on localhost:8001")
	go func() {
		log.Fatalln(s.Serve(lis))
	}()

	// Create a client connection to the gRPC server we just started
	// This is where the gRPC-Gateway proxies the requests
	conn, err := grpc.DialContext(
		context.Background(),
		"0.0.0.0:8001",
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalln("Failed to dial server:", err)
	}

	gwmux := runtime.NewServeMux()
	// Register Greeter
	err = diskpb.RegisterFileServiceHandler(context.Background(), gwmux, conn)
	if err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}

	gwServer := &http.Server{
		Addr:    ":8002",
		Handler: gwmux,
	}

	log.Println("Serving gRPC-Gateway on http://0.0.0.0:8002")
	log.Fatalln(gwServer.ListenAndServe())
}
