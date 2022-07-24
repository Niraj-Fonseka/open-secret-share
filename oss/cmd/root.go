/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	pkg "open-secret-share/oss/pkg"
	"os"
)

// rootCmd represents the base command when called without any subcommands
// var rootCmd = &cobra.Command{
// 	Use:   "oss",
// 	Short: "A brief description of your application",
// 	Long:  `A longer description that spans multiple lines and likely contains`,
// 	// Uncomment the following line if your bare application
// 	// has an action associated with it:
// 	// Run: func(cmd *cobra.Command, args []string) { },
// }

// var initCmd = &cobra.Command{
// 	Use:   "init",
// 	Short: "generate a key pair",
// 	Run:   generateKeyPair,
// }

func Execute() {
	err := pkg.Root.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.oss.yaml)")
	pkg.Root.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	pkg.Root.AddCommand(pkg.Init)
}
