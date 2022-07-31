/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"open-secret-share/oss/client"
	pkg "open-secret-share/oss/pkg"
)

func Execute() {

	client := client.NewKeyServerClient()
	defer client.ConnClose()

	prompt := pkg.NewPrompt()
	pkg.NewCommands(client, prompt).InitializeCommands()

}
