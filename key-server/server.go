package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/ping", GetPubKeyStorage)
	r.GET("/upload", uploadFromMemory)

	r.Run()
}
