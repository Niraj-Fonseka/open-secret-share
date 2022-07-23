package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/alokmenghrajani/gpgeez"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
	"golang.org/x/crypto/openpgp/packet"
)

//helpers https://gist.github.com/ayubmalik/a83ee23c7c700cdce2f8c5bf5f2e9f20

func main() {
	trydecrypt()
}

const pubKey = "pub.gpg"
const fileToEnc = "data.txt"

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
	fmt.Printf("%s\n", output)

	output, err = key.ArmorPrivate(&config)
	if err != nil {
		fmt.Printf("Something went wrong: %v", err)
		return
	}
	fmt.Printf("%s\n", output)

	ioutil.WriteFile("pub.gpg", key.Keyring(), 0666)
	ioutil.WriteFile("priv.gpg", key.Secring(&config), 0666)
}

func trydecrypt() {
	log.Println("Public key:", pubKey)

	// Read in public key
	recipient, err := readEntity(pubKey)
	if err != nil {
		fmt.Println(err)
		return
	}

	f, err := os.Open(fileToEnc)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	dst, err := os.Create(fileToEnc + ".gpg")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer dst.Close()
	encrypt([]*openpgp.Entity{recipient}, nil, f, dst)
}

func encrypt(recip []*openpgp.Entity, signer *openpgp.Entity, r io.Reader, w io.Writer) error {
	wc, err := openpgp.Encrypt(w, recip, signer, &openpgp.FileHints{IsBinary: true}, nil)
	if err != nil {
		return err
	}
	if _, err := io.Copy(wc, r); err != nil {
		return err
	}
	return wc.Close()
}

func readEntity(name string) (*openpgp.Entity, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	block, err := armor.Decode(f)
	if err != nil {
		return nil, err
	}
	return openpgp.ReadEntity(packet.NewReader(block.Body))
}
