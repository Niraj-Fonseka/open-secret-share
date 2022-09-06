package pkg

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"time"
)

type Utils struct {
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func NewUtils() *Utils {
	return &Utils{}
}

func (u *Utils) GenerateUniqueKey() string {
	n := 20
	rand.Seed(time.Now().UnixNano())
	sb := strings.Builder{}
	sb.Grow(n)
	for i := 0; i < n; i++ {
		sb.WriteByte(charset[rand.Intn(len(charset))])
	}
	return sb.String()
}

func (u *Utils) GenerateStorageID(username, uniqueID string) string {
	comb := fmt.Sprintf("%s-%s", username, uniqueID)

	return base64.RawStdEncoding.EncodeToString([]byte(comb))
}

func (u *Utils) GetStorageID(username string) string {
	path := fmt.Sprintf("%s/.oss", os.Getenv("HOME"))

	uniqueID, err := ioutil.ReadFile(path + "/config.txt")

	if err != nil {
		fmt.Printf("unable to read the configuration file. Please re-initialize")
		os.Exit(1)
	}

	return base64.RawStdEncoding.EncodeToString([]byte(fmt.Sprintf("%s-%s", username, string(uniqueID))))
}
