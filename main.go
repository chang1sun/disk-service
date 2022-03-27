package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/changpro/disk-service/application"
	arepo "github.com/changpro/disk-service/domain/auth/repo"
	frepo "github.com/changpro/disk-service/domain/file/repo"
	"github.com/changpro/disk-service/infra/config"
	"github.com/changpro/disk-service/infra/util"
	"github.com/changpro/disk-service/interfaces"
	"github.com/changpro/disk-service/stub"
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
	lis, err := net.Listen("tcp", ":8000")
	if err != nil {
		log.Fatalln("Failed to listen:", err)
	}

	// Create a gRPC server object
	s := grpc.NewServer()
	// Attach the Greeter service to the server
	stub.RegisterAuthServiceServer(s, interfaces.NewAuthServer())
	stub.RegisterFileServiceServer(s, interfaces.NewFileServer())
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

	// independent router for upload interface
	err = AddCustomRoute(gwmux)
	if err != nil {
		panic(err)
	}

	// Register auth RESTful api gateway
	err = stub.RegisterAuthServiceHandler(context.Background(), gwmux, conn)
	if err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}

	// Register file RESTful api gateway
	err = stub.RegisterFileServiceHandler(context.Background(), gwmux, conn)
	if err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}

	gwServer := &http.Server{
		Addr:    ":8001",
		Handler: util.AddMiddleware(gwmux),
	}

	log.Println("Serving gRPC-Gateway on localhost:8001")
	if os.Getenv("RUN_MODE") == "prod" {
		log.Fatalln(gwServer.ListenAndServeTLS(config.GetConfig().TLS.Crt, config.GetConfig().TLS.Key))
	} else {
		log.Fatalln(gwServer.ListenAndServe())
	}
}

// aim to handle file transfer request which cannot be implemented by grpc-gateway
func AddCustomRoute(mux *runtime.ServeMux) error {
	// single small file uplaod
	err := mux.HandlePath("POST", "/v1/file/upload", application.FileUploadHandler)
	if err != nil {
		return err
	}
	// multipart uploader
	err = mux.HandlePath("POST", "/v1/file/mp-upload", application.MPFileUploadHandler)
	if err != nil {
		return err
	}
	// multipart uploader test
	err = mux.HandlePath("POST", "/v1/file/mp-upload-test", application.MPFileUploadTestHandler)
	if err != nil {
		return err
	}
	// download file
	err = mux.HandlePath("GET", "/v1/file/download", application.DownloadHandler)
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
	conn, err := util.GetGormConn(config.GetConfig().Mysql.User, config.GetConfig().Mysql.Password,
		config.GetConfig().Mysql.Addr, config.GetConfig().Mysql.Database)
	if err != nil {
		log.Println(config.GetConfig().Mysql.Addr)
		return err
	}
	mongoDB, err := util.GetMongodbConn(config.GetConfig().MongoDB.Addr, config.GetConfig().MongoDB.Database)
	if err != nil {
		return err
	}
	bucket, err := gridfs.NewBucket(mongoDB)
	if err != nil {
		return err
	}
	redisClient := util.GetRedisConn(config.GetConfig().Redis.Addr,
		config.GetConfig().Redis.User, config.GetConfig().Redis.Password, config.GetConfig().Redis.DBShare)

	// set repo
	arepo.SetUserDao(&arepo.UserDao{Database: conn})
	arepo.SetUserAnalysisDao(&arepo.UserAnalysisDao{Database: conn})
	frepo.SetUniFileStoreDao(&frepo.UniFileStoreDao{Database: mongoDB, Bucket: bucket})
	frepo.SetUserFileDao(&frepo.UserFileDao{Database: mongoDB})
	frepo.SetShareDao(&frepo.ShareDao{Database: redisClient})
	frepo.SetShareRecordDao(&frepo.ShareRecordDao{Database: conn})
	frepo.SetRecycleFileDao(&frepo.RecycleFileDao{Database: mongoDB})
	application.MPRedisClient = util.GetRedisConn(config.GetConfig().Redis.Addr,
		config.GetConfig().Redis.User, config.GetConfig().Redis.Password, config.GetConfig().Redis.DBUpload)
	return nil
}
