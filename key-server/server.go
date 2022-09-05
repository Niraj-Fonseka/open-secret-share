package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"google.golang.org/grpc/codes"

	"google.golang.org/grpc/metadata"

	envconfig "github.com/sethvargo/go-envconfig"

	"open-secret-share/key-server/config"
	cache "open-secret-share/key-server/pkg"
	pb "open-secret-share/key-server/protobuf"
	"open-secret-share/key-server/storageproviders"

	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedOpenSecretShareServer
	Cache   *cache.MemCache
	Storage storageproviders.StorageProvider
}

func (s *server) Recieve(ctx context.Context, in *pb.RecieveRequest) (*pb.RecieveResponse, error) {
	messageID := in.GetMessageId()
	data, found := s.Cache.Get(messageID)

	if !found {
		return &pb.RecieveResponse{}, fmt.Errorf("no message by that id")
	}

	return &pb.RecieveResponse{Data: data}, nil
}

func (s *server) GetPublicKey(ctx context.Context, in *pb.GetPubKeyRequest) (*pb.GetPubKeyResponse, error) {
	username := in.GetUsername()
	data, err := s.Storage.Download(username)

	if err != nil {
		return &pb.GetPubKeyResponse{}, err
	}

	return &pb.GetPubKeyResponse{Pubkey: data}, nil
}

func (s *server) Initialize(ctx context.Context, in *pb.InitializeRequest) (*pb.InitializeResponse, error) {
	pubKey := in.GetPubkey()
	email := in.GetEmail()
	err := s.Storage.Upload(email, pubKey)
	if err != nil {
		return &pb.InitializeResponse{Message: "failed"}, err
	}
	return &pb.InitializeResponse{Message: "success"}, err
}

func (s *server) Store(ctx context.Context, in *pb.StoreRequest) (*pb.StoreResponse, error) {
	data := in.GetEncMessage()

	messageID := s.Cache.Set(data)

	return &pb.StoreResponse{MessageId: messageID}, nil
}

func AuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

	key := os.Getenv("AUTH_KEY")

	if len(key) == 0 {
		fmt.Println("no authentication found ")
		os.Exit(1)
	}

	meta, ok := metadata.FromIncomingContext(ctx)

	if !ok {
		return nil, grpc.Errorf(codes.Unauthenticated, "missing context metadata")
	}
	// Take care: grpc internally reduce key values to lowercase
	if len(meta["authorization"]) != 1 {
		return nil, grpc.Errorf(codes.Unauthenticated, "invalid token")
	}
	if meta["authorization"][0] != key {
		return nil, grpc.Errorf(codes.Unauthenticated, "invalid token")
	}
	return handler(ctx, req)
}

func main() {

	ctx := context.Background()
	var serverConfig config.Server

	if err := envconfig.Process(ctx, &serverConfig); err != nil {
		log.Fatal(err)
	}

	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", serverConfig.PORT))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	cache := cache.NewMemCache()
	storage := storageproviders.NewGoogleStorageClient()
	s := grpc.NewServer(grpc.UnaryInterceptor(AuthInterceptor))
	ossServer := &server{
		Cache:   cache,
		Storage: storage,
	}

	pb.RegisterOpenSecretShareServer(s, ossServer)
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
