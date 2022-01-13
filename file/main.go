package main

import (
	"context"
	"log"
	"net"
	"net/http"

	cutil "github.com/changpro/disk-service/common/util"
	"github.com/changpro/disk-service/file/config"
	"github.com/changpro/disk-service/file/interfaces"
	"github.com/changpro/disk-service/file/repo"
	"github.com/changpro/disk-service/file/service"
	filepb "github.com/changpro/disk-service/pbdeps/file"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Init
	err := InitBase()
	if err != nil {
		panic(err)
	}

	// Create a listener on TCP port
	lis, err := net.Listen("tcp", ":8002")
	if err != nil {
		log.Fatalln("Failed to listen:", err)
	}

	// Create a gRPC server object
	s := grpc.NewServer()
	// Attach the Greeter service to the server
	filepb.RegisterFileServiceServer(s, interfaces.NewServer())
	// Serve gRPC Server
	log.Println("Serving file service on localhost:8002")
	go func() {
		log.Fatalln(s.Serve(lis))
	}()

	// Create a client connection to the gRPC server we just started
	// This is where the gRPC-Gateway proxies the requests
	conn, err := grpc.DialContext(
		context.Background(),
		"localhost:8002",
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalln("Failed to dial server:", err)
	}

	gwmux := runtime.NewServeMux()
	// Register RESTful api gateway
	err = filepb.RegisterFileServiceHandler(context.Background(), gwmux, conn)
	if err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}

	// independent router for upload interface
	err = AddCustomRoute(gwmux)
	if err != nil {
		panic(err)
	}

	gwServer := &http.Server{
		Addr:    ":8003",
		Handler: gwmux,
	}

	log.Println("Serving gRPC-Gateway for file service on localhost:8003")
	log.Fatalln(gwServer.ListenAndServeTLS("crt/service.pem", "crt/service.key"))
}

func AddCustomRoute(mux *runtime.ServeMux) error {
	err := mux.HandlePath("POST", "/v1/files/upload", service.FileUploadHandler)
	if err != nil {
		return err
	}
	err = mux.HandlePath("POST", "v1/files/mp/upload", service.MPFileUploadHandler)
	if err != nil {
		return err
	}
	err = mux.HandlePath("POST", "v1/files/mp/upload-finish", service.FileMergeHandler)
	if err != nil {
		return err
	}
	return nil
}

func InitBase() error {
	err := config.InitConfig()
	if err != nil {
		return err
	}
	mongoDB, err := cutil.GetMongodbConn(config.GetConfig().MongoDB.Addr, config.GetConfig().MongoDB.Database)
	if err != nil {
		return err
	}
	bucket, err := gridfs.NewBucket(mongoDB)
	if err != nil {
		return err
	}
	repo.SetUniFileStoreDao(&repo.UniFileStoreDao{Database: mongoDB, Bucket: bucket})
	return nil
}
