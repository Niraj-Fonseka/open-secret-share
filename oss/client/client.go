package client

import (
	"flag"

	"google.golang.org/grpc"

	"google.golang.org/grpc/credentials/insecure"

	"log"
	pb "open-secret-share/oss/protobuf"
)

type KeyServerClient struct {
	Client pb.GreeterClient
	conn   *grpc.ClientConn
}

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
)

func NewKeyServerClient() *KeyServerClient {
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	c := pb.NewGreeterClient(conn)

	return &KeyServerClient{Client: c, conn: conn}
}

func (k *KeyServerClient) ConnClose() {
	k.conn.Close()
}
