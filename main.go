package main

import (
	"log"
	"net"
	"os"

	pb "github.com/jimdhughes/go-user-ms/proto"
	"google.golang.org/grpc"
)

func main() {
	serverType := os.Getenv("USERMS_SERVER_TYPE")
	if serverType == "" {
		serverType = "http"
	}
	InitializeDatabase()
	InitializeTokenService()
	if serverType == "grpc" {
		InitializeGRPCService()
	}
	if serverType == "http" {
		InitializeHttpServer()
	}
}

func InitializeGRPCService() {
	log.Println("initializing grpc service")
	grpcPort := os.Getenv("USERMS_GRPC_PORT")
	if grpcPort == "" {
		grpcPort = ":50051"
	}

	lis, err := net.Listen("tcp", grpcPort)
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
	log.Println("initializing database")
	DB = &DBClient{}
	DB.Initialize("./data/bolt.db")
}

func InitializeTokenService() {
	log.Println("initializing token service")
	TS = &TokenService{}
}
