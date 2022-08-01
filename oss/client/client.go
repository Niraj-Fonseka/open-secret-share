package client

import (
	"flag"
	"os"

	"google.golang.org/grpc"

	"google.golang.org/grpc/credentials/insecure"

	"log"
	pb "open-secret-share/oss/protobuf"
)

type KeyServerClient struct {
	Client pb.OpenSecretShareClient
	conn   *grpc.ClientConn
}

func NewKeyServerClient() *KeyServerClient {

	serverEndpoint := os.Getenv("SERVER")

	if serverEndpoint == "" {
		log.Fatal("SERVER environment variable is not set")
	}

	addr := flag.String("addr", serverEndpoint, "key server address")

	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	c := pb.NewOpenSecretShareClient(conn)

	return &KeyServerClient{Client: c, conn: conn}
}

func (k *KeyServerClient) ConnClose() {
	k.conn.Close()
}
