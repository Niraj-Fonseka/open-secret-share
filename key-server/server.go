package main

import (
	"context"
	"log"
	"net/http"

	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/ping", GetPubKey)

	r.Run()
}

func GetPubKey(c *gin.Context) {
	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		// TODO: Handle error.
		log.Fatal(err)

	}
	defer client.Close()

	req := &secretmanagerpb.CreateSecretRequest{
		Parent: "keyname",
		Secret: &secretmanagerpb.Secret{Name: "keyname"},
		// TODO: Fill request struct fields.
		// See https://pkg.go.dev/google.golang.org/genproto/googleapis/cloud/secretmanager/v1#CreateSecretRequest.
	}
	resp, err := client.CreateSecret(ctx, req)
	if err != nil {
		// TODO: Handle error.
		log.Fatal(err)
	}
	// TODO: Use resp.
	_ = resp
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
