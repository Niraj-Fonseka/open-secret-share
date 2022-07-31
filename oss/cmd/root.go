/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	pkg "open-secret-share/oss/pkg"
)

func Execute() {

	pkg.NewCommands().InitializeCommands()

}
