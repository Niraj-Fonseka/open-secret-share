package pkg

import (
	"context"
	"fmt"
	"log"
	"open-secret-share/oss/client"
	pb "open-secret-share/oss/protobuf"
	"os"
	"strings"
	"time"

	"google.golang.org/grpc/metadata"

	"github.com/spf13/cobra"
)

type Commands struct {
	Root        *cobra.Command
	Send        *cobra.Command
	Init        *cobra.Command
	RecieveCMD  *cobra.Command
	UsernameCMD *cobra.Command
	prompt      *Prompt
	client      pb.OpenSecretShareClient
	gpgTools    *GPGTools
	utils       *Utils
}

//NewCommands
//Create a Commands Client
func NewCommands(client *client.KeyServerClient, prompt *Prompt, gpgtools *GPGTools, utils *Utils) *Commands {
	return &Commands{prompt: prompt, client: client.Client, gpgTools: gpgtools, utils: utils}
}

//InitializeCommands
//Initialize all the commands
func (c *Commands) InitializeCommands() *Commands {
	var commands Commands

	//root command
	commands.Root = &cobra.Command{
		Use:   "oss",
		Short: "app for sharing secrets",
		Long:  `A longer description that spans multiple lines and likely contains`,
		Run:   c.defaultHandler,
	}

	//command for initilaizing
	commands.Init = &cobra.Command{
		Use:   "init",
		Short: "generate a new key pair and initialize the app",
		Run:   c.initializeHandler,
	}

	//command for sending
	commands.Send = &cobra.Command{
		Use:   "send",
		Short: "send a message to a user",
		Run:   c.sendSecretHandler,
	}

	//command for receiving
	commands.RecieveCMD = &cobra.Command{
		Use:   "recieve",
		Short: "receive a message given id",
		Run:   c.recieveHandler,
	}

	commands.UsernameCMD = &cobra.Command{
		Use:   "username",
		Short: "get username",
		Run:   c.usernameHandler,
	}

	commands.Root.AddCommand(commands.Init)
	commands.Root.AddCommand(commands.Send)
	commands.Root.AddCommand(commands.RecieveCMD)
	commands.Root.AddCommand(commands.UsernameCMD)

	err := commands.Root.Execute()
	if err != nil {
		os.Exit(1)
	}
	return &commands
}

/*
	initializeHandler
	- get user input for username, email and comment and generate a pub/pvt key pair
	- store the pub key in the key server and store the private key locally
*/
func (c *Commands) initializeHandler(cmd *cobra.Command, args []string) {
	username := c.prompt.TriggerPrompt("username")
	email := c.prompt.TriggerPrompt("email")
	comment := c.prompt.TriggerPrompt("comment")

	sanitizedUN := c.gpgTools.utils.SanitizeUsername(username)
	pubKey := c.gpgTools.GenerateKeyPair(sanitizedUN, strings.TrimSpace(email), strings.TrimSpace(comment))

	key := os.Getenv("AUTH_KEY")

	if len(key) == 0 {
		fmt.Println("no authentication found ")
		os.Exit(1)
	}

	md := metadata.Pairs("Authorization", key)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(metadata.NewOutgoingContext(context.Background(), md), time.Second*120)
	defer cancel()

	//peek first

	r, err := c.client.Initialize(ctx, &pb.InitializeRequest{Pubkey: pubKey, Username: c.gpgTools.utils.SanitizeUsername(username)})
	if err != nil {
		fmt.Printf("Error happened at initialization : %s", err.Error())
		os.Exit(2)
	}

	err = c.utils.WriteConfig(sanitizedUN)
	if err != nil {
		fmt.Printf("Error happened at initialization : %s", err.Error())
		os.Exit(3)
	}

	fmt.Println(r.Message)
}

/*
	sendSecretHandler
	- prompt the user for the reciever's email
	- recieve the public key from the key server
	- prompt the user for the message to be sent
	- encrpyt the message using reciever's pub key
	- store in memory cache. generate a unique id
*/
func (c *Commands) sendSecretHandler(cmd *cobra.Command, args []string) {

	c.utils.CheckInitialized()
	receiver := c.prompt.TriggerPrompt("reciever username")
	username := c.utils.SanitizeUsername(receiver)

	key := os.Getenv("AUTH_KEY")

	if len(key) == 0 {
		fmt.Println("no authentication found ")
		os.Exit(1)
	}

	md := metadata.Pairs("Authorization", key)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(metadata.NewOutgoingContext(context.Background(), md), time.Second*120)
	defer cancel()

	fmt.Print("----- establishing connection with the key server ----- ")
	r, err := c.client.GetPublicKey(ctx, &pb.GetPubKeyRequest{Username: username})
	if err != nil {
		fmt.Printf("unable to establish connection with the server %v", err)

		os.Exit(1)
	}

	pubKeyRecived := r.GetPubkey() //probably can return an error also. When the pub key isnt there

	if len(pubKeyRecived) == 0 {
		fmt.Printf("unable to find the public key of %s \n", receiver)
		return
	}

	fmt.Println("connection established successfully")
	message := c.prompt.TriggerPrompt("message")

	encrypted, err := c.gpgTools.Encrypt(message, pubKeyRecived)

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

/*
	recieveHandler
	- prompt the user for the message id
	- receive the encrypted message from the key server
	- decrypt the data using the reciever's private key
*/
func (c *Commands) recieveHandler(cmd *cobra.Command, args []string) {

	c.utils.CheckInitialized()

	messageID := c.prompt.TriggerPrompt("message id")

	key := os.Getenv("AUTH_KEY")

	if len(key) == 0 {
		fmt.Println("no authentication found ")
		os.Exit(1)
	}

	md := metadata.Pairs("Authorization", key)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(metadata.NewOutgoingContext(context.Background(), md), time.Second*120)
	defer cancel()

	r, err := c.client.Recieve(ctx, &pb.RecieveRequest{MessageId: messageID})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	encData := r.GetData()

	decrypted, err := c.gpgTools.Decrypt(encData)

	fmt.Println(decrypted)

}

func (c *Commands) usernameHandler(cmd *cobra.Command, args []string) {

	c.utils.CheckInitialized()

	config, err := c.utils.ReadConfig()

	if err != nil {
		fmt.Println("no config found. Please initialize the app")
		return
	}

	fmt.Println(config.Username)
}

/*
	defaultHandler
	- list out the available commands
*/
func (c *Commands) defaultHandler(cmd *cobra.Command, args []string) {
	fmt.Println("Hello from Open Secret Share !")
	fmt.Println("oss init - to initialize the app")
	fmt.Println("oss send - to send a message")
	fmt.Println("oss recieve - to recieve a message")
	fmt.Println("oss recieve - to recieve a message")
	fmt.Println("oss username - fetch the username")

}
