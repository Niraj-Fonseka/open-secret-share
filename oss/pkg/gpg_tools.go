package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/alokmenghrajani/gpgeez"
)

func GenerateKeyPair() []byte {
	config := gpgeez.Config{Expiry: 365 * 24 * time.Hour}
	key, err := gpgeez.CreateKey("JoeJoe", "test key", "joe@example.com", &config)
	if err != nil {
		fmt.Printf("Something went wrong: %v", err)
		return []byte{}
	}
	_, err = key.Armor()
	if err != nil {
		fmt.Printf("Something went wrong: %v", err)
		return []byte{}
	}

	_, err = key.ArmorPrivate(&config)
	if err != nil {
		fmt.Printf("Something went wrong: %v", err)
		return []byte{}
	}

	pub := key.Keyring()
	pvt := key.Secring(&config)
	pub_err := ioutil.WriteFile("oss_pub.gpg", pub, 0666)
	if err != nil {
		log.Printf("error when writing public key : %v\n", pub_err)
		return []byte{}
	}
	pvt_err := ioutil.WriteFile("oss_pvt.gpg", pvt, 0666)
	if err != nil {
		log.Printf("error when writing public key : %v\n", pvt_err)
		return []byte{}
	}

	return pub
}
