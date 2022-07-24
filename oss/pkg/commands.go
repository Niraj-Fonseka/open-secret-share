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
	// Run: func(cmd *cobra.Command, args []string) { },
}

var Init = &cobra.Command{
	Use:   "init",
	Short: "generate a key pair",
	Run:   InitializeApp,
}
