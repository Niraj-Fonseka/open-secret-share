package pkg

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"time"
)

type Utils struct {
}

type Config struct {
	Username string `json:"username"`
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

func (u *Utils) SanitizeUsername(username string) string {
	return strings.ToLower(strings.Replace(username, " ", "_", -1)) //repalce all the spaces with underscores
}

func (u *Utils) WriteConfig(username string) error {
	path := fmt.Sprintf("%s/.oss", os.Getenv("HOME"))
	var config Config
	config.Username = username

	bts, err := json.Marshal(&config)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path+"/config.json", bts, 0666)
}

func (u *Utils) ReadConfig() (Config, error) {
	var config Config
	path := fmt.Sprintf("%s/.oss", os.Getenv("HOME"))
	configFile, err := os.ReadFile(path + "/config.json")
	if err != nil {
		return config, err
	}

	err = json.Unmarshal(configFile, &config)

	if err != nil {
		return config, err
	}

	return config, err
}

func (u *Utils) CheckInitialized() {
	path := fmt.Sprintf("%s/.oss", os.Getenv("HOME"))

	_, err := os.Stat(path)

	if os.IsNotExist(err) {
		fmt.Println("can't find the init directory please initialize your app")
		os.Exit(1)
	}

	_, err = os.Stat(fmt.Sprintf("%s/%s", path, "config.json"))

	if os.IsNotExist(err) {
		fmt.Println("invalid configuration. Please re-initialize")
		os.Exit(2)
	}

	_, err = os.Stat(fmt.Sprintf("%s/%s", path, "oss_pub.gpg"))

	if os.IsNotExist(err) {
		fmt.Println("invalid configuration. Please re-initialize")
		os.Exit(2)
	}

	_, err = os.Stat(fmt.Sprintf("%s/%s", path, "oss_pvt.gpg"))

	if os.IsNotExist(err) {
		fmt.Println("invalid configuration. Please re-initialize")
		os.Exit(3)
	}

}
