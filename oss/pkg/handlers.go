package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

func InitializeApp(cmd *cobra.Command, args []string) {
	log.Println("Generate Key Pair...")
	GenerateKeyPair()
}
