package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/changpro/disk-service/auth/config"
	"github.com/changpro/disk-service/auth/interfaces"
	"github.com/changpro/disk-service/auth/repo"
	cutil "github.com/changpro/disk-service/common/util"
	authpb "github.com/changpro/disk-service/pbdeps/auth"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// handle cors
func cors(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.Header.Get("Origin"), "http://localhost") {
			w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization, ResponseType")
		}
		if r.Method == "OPTIONS" {
			return
		}
		h.ServeHTTP(w, r)
	})
}

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
		runtime.WithErrorHandler(cutil.CustomErrorHandler),
	)
	// Register RESTful api gateway
	err = authpb.RegisterAuthServiceHandler(context.Background(), gwmux, conn)
	if err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}

	gwServer := &http.Server{
		Addr:    ":8001",
		Handler: cors(gwmux),
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
	err := config.InitConfig()
	if err != nil {
		return err
	}
	conn, err := cutil.GetGormConn(config.GetConfig().Mysql.User, config.GetConfig().Mysql.Password,
		config.GetConfig().Mysql.Addr, config.GetConfig().Mysql.Database)
	if err != nil {
		return err
	}
	repo.SetUserDao(&repo.UserDao{Database: conn})
	repo.SetUserAnalysisDao(&repo.UserAnalysisDao{Database: conn})
	return nil
}
