/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	pkg "open-secret-share/oss/pkg"
	"os"
)

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
