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

func GenerateKeyPair(username, email, comment string) []byte {
	config := gpgeez.Config{Expiry: 365 * 24 * time.Hour}
	key, err := gpgeez.CreateKey(username, comment, email, &config)
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

func Encrypt(data string, pubKey []byte) (string, error) {
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

	// Output encrypted/encoded string
	log.Println("Encrypted Secret:", encStr)

	return encStr, nil
}

func Decrypt(encryptedString string) (string, error) {

	const passphrase = "" //go/crypto doesn't support passpharse yet

	// init some vars
	var entity *openpgp.Entity
	var entityList openpgp.EntityList

	// Open the private key file
	keyringFileBuffer, err := os.Open("/home/hungryotter/go/src/open-secret-share/oss/oss_pvt.gpg")
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
	entity.PrivateKey.Decrypt(passphraseByte)
	for _, subkey := range entity.Subkeys {
		subkey.PrivateKey.Decrypt(passphraseByte)
	}
	log.Println("Finished decrypting private key using passphrase")

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
