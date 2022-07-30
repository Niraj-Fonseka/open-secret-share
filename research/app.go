package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/alokmenghrajani/gpgeez"
	"golang.org/x/crypto/openpgp"
)

/*
helpers
https://gist.github.com/ayubmalik/a83ee23c7c700cdce2f8c5bf5f2e9f20 <- using asc
https://gist.github.com/stuart-warren/93750a142d3de4e8fdd2 <- using gpg ( we using this now )
*/
const mySecretString = "this is so very secret!"
const passphrase = "test asdkey"
const secretKeyring = "priv.gpg"
const publicKeyring = "pub2.gpg"

func main() {
	encStr, err := encTest(mySecretString)
	if err != nil {
		log.Fatal(err)
	}
	decStr, err := decTest(encStr)
	if err != nil {
		log.Fatal(err)
	}
	// should be done
	log.Println("Decrypted Secret:", decStr)
}

//need to take in a user data structure
func generatekeys() {
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

	ioutil.WriteFile("pub2.gpg", key.Keyring(), 0666)
	ioutil.WriteFile("priv2.gpg", key.Secring(&config), 0666)
}

func encTest(secretString string) (string, error) {
	log.Println("Secret to hide:", secretString)
	log.Println("Public Keyring:", publicKeyring)

	// Read in public key
	keyringFileBuffer, _ := os.Open(publicKeyring)
	defer keyringFileBuffer.Close()
	entityList, err := openpgp.ReadKeyRing(keyringFileBuffer)
	if err != nil {
		return "", err
	}

	// encrypt string
	buf := new(bytes.Buffer)
	w, err := openpgp.Encrypt(buf, entityList, nil, nil, nil)
	if err != nil {
		return "", err
	}
	_, err = w.Write([]byte(mySecretString))
	if err != nil {
		return "", err
	}
	err = w.Close()
	if err != nil {
		return "", err
	}

	// Encode to base64
	bytes, err := ioutil.ReadAll(buf)
	if err != nil {
		return "", err
	}
	encStr := base64.StdEncoding.EncodeToString(bytes)

	// Output encrypted/encoded string
	log.Println("Encrypted Secret:", encStr)

	return encStr, nil
}

func decTest(encString string) (string, error) {

	log.Println("Secret Keyring:", secretKeyring)
	log.Println("Passphrase:", passphrase)

	// init some vars
	var entity *openpgp.Entity
	var entityList openpgp.EntityList

	// Open the private key file
	keyringFileBuffer, err := os.Open(secretKeyring)
	if err != nil {
		return "", err
	}
	defer keyringFileBuffer.Close()
	entityList, err = openpgp.ReadKeyRing(keyringFileBuffer)
	if err != nil {
		return "", err
	}
	entity = entityList[0]

	// Get the passphrase and read the private key.
	// Have not touched the encrypted string yet
	passphraseByte := []byte(passphrase)
	log.Println("Decrypting private key using passphrase")
	entity.PrivateKey.Decrypt(passphraseByte)
	for _, subkey := range entity.Subkeys {
		subkey.PrivateKey.Decrypt(passphraseByte)
	}
	log.Println("Finished decrypting private key using passphrase")

	// Decode the base64 string
	dec, err := base64.StdEncoding.DecodeString(encString)
	if err != nil {
		return "", err
	}

	// Decrypt it with the contents of the private key
	md, err := openpgp.ReadMessage(bytes.NewBuffer(dec), entityList, nil, nil)
	if err != nil {
		return "", err
	}
	bytes, err := ioutil.ReadAll(md.UnverifiedBody)
	if err != nil {
		return "", err
	}
	decStr := string(bytes)

	return decStr, nil
}
