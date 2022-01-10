package main

import (
	"context"
	"log"
	"net"
	"net/http"

	authpb "github.com/changpro/disk-service/auth/deps"
	"github.com/changpro/disk-service/auth/interfaces"
	"github.com/changpro/disk-service/auth/repo"
	"github.com/changpro/disk-service/auth/util"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Init
	err := InitBase()
	if err != nil {
		panic(err)
	}
	InitDaoImpl()

	// Create a listener on TCP port
	lis, err := net.Listen("tcp", ":8000")
	if err != nil {
		log.Fatalln("Failed to listen:", err)
	}

	// Create a gRPC server object
	s := grpc.NewServer()
	// Attach the Greeter service to the server
	authpb.RegisterAuthServiceServer(s, interfaces.NewServer())
	// Serve gRPC Server
	log.Println("Serving gRPC on localhost:8000")
	go func() {
		err := s.Serve(lis)
		if err != nil {
			log.Fatalln(err)
			panic(err)
		}
	}()
	// Create a client connection to the gRPC server we just started
	// This is where the gRPC-Gateway proxies the requests
	conn, err := grpc.DialContext(
		context.Background(),
		"localhost:8000",
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalln("Failed to dial server:", err)
	}
	gwmux := runtime.NewServeMux(
		runtime.WithErrorHandler(util.CustomErrorHandler),
	)
	// Register RESTful api gateway
	err = authpb.RegisterAuthServiceHandler(context.Background(), gwmux, conn)
	if err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}

	gwServer := &http.Server{
		Addr:    ":8001",
		Handler: gwmux,
	}

	log.Println("Serving gRPC-Gateway on localhost:8001")
	// err = gwServer.ListenAndServeTLS("crt/service.pem", "crt/service.key")
	err = gwServer.ListenAndServe()
	if err != nil {
		log.Fatalln(err)
		panic(err)
	}
}

func InitBase() error {
	err := util.InitConfig()
	if err != nil {
		return err
	}
	err = util.InitGormConn()
	if err != nil {
		return err
	}
	return nil
}

func InitDaoImpl() {
	repo.SetUserDao(&repo.UserDao{Database: util.GetGormConn()})
	repo.SetUserAnalysisDao(&repo.UserAnalysisDao{Database: util.GetGormConn()})
}
