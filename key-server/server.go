package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	storage "cloud.google.com/go/storage"
	"google.golang.org/api/iterator"

	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/ping", GetPubKeySM)

	r.Run()
}

func GetPubKeySrorage(c *gin.Context) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		// TODO: Handle error.
		log.Fatal(err)
	}

	bkt := client.Bucket("oss-pubkeys")

	query := &storage.Query{Prefix: ""}

	var names []string
	it := bkt.Objects(ctx, query)
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		names = append(names, attrs.Name)
	}
	fmt.Println(names)
}

func GetPubKeySM(c *gin.Context) {
	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		// TODO: Handle error.
		log.Fatal(err)

	}
	defer client.Close()

	req := &secretmanagerpb.CreateSecretRequest{
		Parent: "test-project-1223",
		Secret: &secretmanagerpb.Secret{Name: "keyname"},
		// TODO: Fill request struct fields.
		// See https://pkg.go.dev/google.golang.org/genproto/googleapis/cloud/secretmanager/v1#CreateSecretRequest.
	}
	resp, err := client.CreateSecret(ctx, req)
	if err != nil {
		// TODO: Handle error.
		log.Fatal(err)
	}

	fmt.Println(resp)
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
