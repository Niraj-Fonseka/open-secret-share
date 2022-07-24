package cmd

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/alokmenghrajani/gpgeez"
)

func GenerateKeyPair() {
	config := gpgeez.Config{Expiry: 365 * 24 * time.Hour}
	key, err := gpgeez.CreateKey("JoeJoe", "test key", "joe@example.com", &config)
	if err != nil {
		fmt.Printf("Something went wrong: %v", err)
		return
	}
	output, err := key.Armor()
	if err != nil {
		fmt.Printf("Something went wrong: %v", err)
		return
	}
	fmt.Println("getting here")
	fmt.Printf("%s\n", output)

	output, err = key.ArmorPrivate(&config)
	if err != nil {
		fmt.Printf("Something went wrong: %v", err)
		return
	}

	fmt.Println("getting here")
	fmt.Printf("%s\n", output)

	ioutil.WriteFile("pub.gpg", key.Keyring(), 0666)
	ioutil.WriteFile("pvt.gpg", key.Secring(&config), 0666)
}
