package cmd

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	pb "open-secret-share/oss/protobuf"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/spf13/cobra"
)

const (
	defaultName = "world"
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
	name = flag.String("name", defaultName, "Name to greet")
)

func InitializeApp(cmd *cobra.Command, args []string) {
	log.Println("Initializing app")
	GenerateKeyPair()
}

func SendHandler(cmd *cobra.Command, args []string) {
	log.Println("sending to ")
	username := "fonseka_live"
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.Recieve(ctx, &pb.RecieveRequest{Username: username})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	pubKeyRecived := r.GetPubkey()

	fmt.Println("pub key : ", pubKeyRecived)
	//encrypt the data file
	//upload the encryptd data into memory
	//generate a uniuq  indentifier============================================================================================
}

func Test(cmd *cobra.Command, args []string) {
	flag.Parse()
	// Set up a connection to the server.
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)
	pubkey := GenerateKeyPair()

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.Initialize(ctx, &pb.InitializeRequest{Pubkey: pubkey})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.GetMessage())
}
