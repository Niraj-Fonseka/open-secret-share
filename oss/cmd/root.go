/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"open-secret-share/oss/client"
	pkg "open-secret-share/oss/pkg"
)

//Execute
//Initialize and Execute all the commands
func Execute() {

	client := client.NewKeyServerClient()
	defer client.ConnClose()

	prompt := pkg.NewPrompt()
	gpgTools := pkg.NewGPGTools()
	pkg.NewCommands(client, prompt, gpgTools).InitializeCommands()

}
