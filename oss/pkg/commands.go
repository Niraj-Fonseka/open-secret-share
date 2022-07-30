package cmd

import (
	"github.com/spf13/cobra"
)

var Root = &cobra.Command{
	Use:   "oss",
	Short: "A brief description of your application",
	Long:  `A longer description that spans multiple lines and likely contains`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: Test,
}

var Init = &cobra.Command{
	Use:   "init",
	Short: "generate a new key pair and initialize the app",
	Run:   InitializeApp,
}

var Send = &cobra.Command{
	Use:   "send",
	Short: "send a message to a user",
	Run:   SendSecret,
}

var RecieveCMD = &cobra.Command{
	Use:   "recieve",
	Short: "receive a message given id",
	Run:   Recieve,
}
