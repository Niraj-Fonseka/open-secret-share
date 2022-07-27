package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"open-secret-share/key-server/pkg"
	pb "open-secret-share/key-server/protobuf"

	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedGreeterServer
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v", in.GetName())
	storage := pkg.NewStorageClient()
	storage.Upload("fonseka_live_gmail", []byte(in.GetName())) //replace all special characters with underscores
	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func (s *server) Initialize(ctx context.Context, in *pb.InitializeRequest) (*pb.InitializeResponse, error) {
	pubKey := in.GetPubkey()
	storage := pkg.NewStorageClient()
	err := storage.Upload("fonseka_live_gmail", pubKey)
	if err != nil {
		return &pb.InitializeResponse{}, err
	}
	return &pb.InitializeResponse{Message: "success"}, err
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
