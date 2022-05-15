package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	pb "github.com/jimdhughes/go-user-ms/proto"
	"google.golang.org/grpc"
)

func main() {
	serverType := flag.String("serverType", "http", "start up an http or grpc server")
	flag.Parse()
	fmt.Println("server type:", *serverType)
	InitializeDatabase()
	InitializeTokenService()
	if *serverType == "grpc" {
		InitializeGRPCService()
	}
	if *serverType == "http" {
		InitializeHttpServer()
	}
}

func InitializeGRPCService() {
	log.Println("initializing grpc service")
	lis, err := net.Listen("tcp", ":3031")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterUserServiceServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func InitializeHttpServer() {
	log.Println("initializing http server")
	InitRouter()
}

func InitializeDatabase() {
	log.Println("initializing databaase")
	DB = &DBClient{}
	DB.Initialize("./data/bolt.db")
}

func InitializeTokenService() {
	log.Println("initializing token service")
	TS = &TokenService{}
}
