package pkg

import (
	"context"
	"flag"
	"fmt"
	"log"
	pb "open-secret-share/oss/protobuf"
	"os"
	"time"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Commands struct {
	Root       *cobra.Command
	Send       *cobra.Command
	Init       *cobra.Command
	RecieveCMD *cobra.Command
	prompt     *Prompt
	client     pb.GreeterClient
}

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
)

func NewCommands() *Commands {
	prompt := NewPrompt()

	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	return &Commands{prompt: prompt, client: c}
}

func (c *Commands) InitializeCommands() *Commands {
	var commands Commands

	//List of commands
	commands.Root = &cobra.Command{
		Use:   "oss",
		Short: "app for sharing secrets",
		Long:  `A longer description that spans multiple lines and likely contains`,
		Run:   c.defaultHandler,
	}

	commands.Init = &cobra.Command{
		Use:   "init",
		Short: "generate a new key pair and initialize the app",
		Run:   c.initializeHandler,
	}

	commands.Send = &cobra.Command{
		Use:   "send",
		Short: "send a message to a user",
		Run:   c.sendSecretHandler,
	}

	commands.RecieveCMD = &cobra.Command{
		Use:   "recieve",
		Short: "receive a message given id",
		Run:   c.recieveHandler,
	}

	commands.Root.AddCommand(commands.Init)
	commands.Root.AddCommand(commands.Send)
	commands.Root.AddCommand(commands.RecieveCMD)

	err := commands.Root.Execute()
	if err != nil {
		os.Exit(1)
	}
	return &commands
}

func (c *Commands) initializeHandler(cmd *cobra.Command, args []string) {
	log.Println("Initializing app")
	username := c.prompt.TriggerPrompt("username :")
	email := c.prompt.TriggerPrompt("email :")
	comment := c.prompt.TriggerPrompt("comment :")

	pubKey := GenerateKeyPair(username, email, comment)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*120)
	defer cancel()

	r, err := c.client.Initialize(ctx, &pb.InitializeRequest{Pubkey: pubKey, Email: email})
	if err != nil {
		log.Fatalf("Error happened at initialization.  : %v", err)
	}

	fmt.Println(r.Message)
}

func (c *Commands) sendSecretHandler(cmd *cobra.Command, args []string) {
	log.Println("preparing to send a secret")

	receiver := c.prompt.TriggerPrompt("email :")
	username := receiver

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	log.Println("----- sending a request to the key-server ...")
	fmt.Println(c.client)
	r, err := c.client.GetPublicKey(ctx, &pb.GetPubKeyRequest{Username: username})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	pubKeyRecived := r.GetPubkey() //probably can return an error also. When the pub key isnt there

	if len(pubKeyRecived) == 0 {
		log.Printf("unable to find the public key of %s \n", receiver)
		return
	}

	fmt.Println("Public Key Recieved ..")
	message := c.prompt.TriggerPrompt("message :")

	encrypted, err := Encrypt(message, pubKeyRecived)

	if err != nil {
		log.Println(err)
		return
	}

	messageID, err := c.client.Store(ctx, &pb.StoreRequest{EncMessage: encrypted})

	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println(messageID.MessageId)

}

func (c *Commands) recieveHandler(cmd *cobra.Command, args []string) {
	messageID := c.prompt.TriggerPrompt("message id :")

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	r, err := c.client.Recieve(ctx, &pb.RecieveRequest{MessageId: messageID})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	encData := r.GetData()

	decrypted, err := Decrypt(encData)

	fmt.Println("Decrupted Message : ", decrypted)

}

func (c *Commands) defaultHandler(cmd *cobra.Command, args []string) {
	fmt.Println("Hello from Open Secret Share !")
	fmt.Println("run oss init - to initialize the app")
	fmt.Println("run oss send - to send a message")
	fmt.Println("run oss recieve - to recieve a message")
}
