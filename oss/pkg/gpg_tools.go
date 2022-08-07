package pkg

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"golang.org/x/crypto/openpgp"

	"github.com/alokmenghrajani/gpgeez"
)

type GPGTools struct {
}

//NewGPGTools
//Create a new GPGTools client
func NewGPGTools() *GPGTools {
	return &GPGTools{}
}

/*
	GenerateKeyPair
	- generate gpg key pair
*/
func (g *GPGTools) GenerateKeyPair(username, email, comment string) []byte {
	path := fmt.Sprintf("%s/.oss", os.Getenv("HOME"))
	config := gpgeez.Config{Expiry: 365 * 24 * time.Hour}
	key, err := gpgeez.CreateKey(username, comment, email, &config)
	if err != nil {
		fmt.Printf("Something went wrong: %v", err)
		os.Exit(1)
		return []byte{}
	}
	_, err = key.Armor()
	if err != nil {
		fmt.Printf("Something went wrong: %v", err)
		os.Exit(1)
		return []byte{}
	}

	_, err = key.ArmorPrivate(&config)
	if err != nil {
		fmt.Printf("Something went wrong: %v", err)
		os.Exit(1)
		return []byte{}
	}

	_, err = os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			err := os.MkdirAll(path, 0755)
			if err != nil {
				fmt.Printf("unable to create the required directory at path : %s", path)
				os.Exit(1)
			}
		} else {
			fmt.Printf("something went wrong when creating the required directory  : %v", err)
			os.Exit(1)
		}
	}

	pub := key.Keyring()
	pvt := key.Secring(&config)
	pub_err := ioutil.WriteFile(path+"/oss_pub.gpg", pub, 0666)
	if pub_err != nil {
		log.Printf("error when writing public key : %v\n", pub_err)
		os.Exit(1)

		return []byte{}
	}
	pvt_err := ioutil.WriteFile(path+"/oss_pvt.gpg", pvt, 0666)
	if pvt_err != nil {
		log.Printf("error when writing the private key : %v\n", pvt_err)
		os.Exit(1)

		return []byte{}
	}

	return pub
}

/*
	Encrypt
	- encrypt a string of data with a public key
*/
func (g *GPGTools) Encrypt(data string, pubKey []byte) (string, error) {
	// Read in public key
	publicKeyring := bytes.NewReader(pubKey)

	entityList, err := openpgp.ReadKeyRing(publicKeyring)
	if err != nil {
		return "", err
	}

	// encrypt string
	buf := new(bytes.Buffer)
	w, err := openpgp.Encrypt(buf, entityList, nil, nil, nil)
	if err != nil {
		return "", err
	}
	_, err = w.Write([]byte(data))
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

	return encStr, nil
}

/*
	Decrypt
	- decrpyt an encrypted string using the private key
*/
func (g *GPGTools) Decrypt(encryptedString string) (string, error) {
	path := fmt.Sprintf("%s/.oss", os.Getenv("HOME"))

	const passphrase = "" //go/crypto doesn't support passpharse yet

	// init some vars
	var entity *openpgp.Entity
	var entityList openpgp.EntityList

	_, err := os.Stat(path)
	if err != nil {
		fmt.Printf("uninitialized. please run the command -> oss init")
		os.Exit(1)
	}
	if os.IsNotExist(err) {
		fmt.Printf("uninitialized. please run the command -> oss init")
		os.Exit(1)
	}

	// Open the private key file
	keyringFileBuffer, err := os.Open(path + "/oss_pvt.gpg")
	if err != nil {
		log.Println("unable to find the private key : please re-initialize the app")
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
	entity.PrivateKey.Decrypt(passphraseByte)
	for _, subkey := range entity.Subkeys {
		subkey.PrivateKey.Decrypt(passphraseByte)
	}

	// Decode the base64 string
	dec, err := base64.StdEncoding.DecodeString(encryptedString)
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
