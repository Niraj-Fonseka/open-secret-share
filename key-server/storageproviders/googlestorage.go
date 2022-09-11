package storageproviders

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"log"
	"open-secret-share/key-server/config"
	"time"

	envconfig "github.com/sethvargo/go-envconfig"

	"google.golang.org/api/option"

	storage "cloud.google.com/go/storage"
)

type GoogleStorage struct {
	ctx    context.Context
	client *storage.Client
	bkt    *storage.BucketHandle
}

//Create a new google storage client
func NewGoogleStorageClient() *GoogleStorage {
	ctx := context.Background()

	var googleConfig config.GoogleStorage

	if err := envconfig.Process(ctx, &googleConfig); err != nil {
		log.Fatal(err)
	}

	client, err := storage.NewClient(ctx, option.WithCredentialsJSON([]byte(googleConfig.GoogleServiceAccount)))

	if err != nil {
		log.Fatalf("Unable to initialize the storage client %v ", err)
	}

	bkt := client.Bucket(googleConfig.BucketName)

	return &GoogleStorage{
		ctx:    ctx,
		client: client,
		bkt:    bkt,
	}
}

func (s *GoogleStorage) Close() {
	s.client.Close()
}

func (s *GoogleStorage) Upload(userID string, pubkey []byte) error {

	buf := bytes.NewBuffer(pubkey)

	ctx, cancel := context.WithTimeout(s.ctx, time.Second*120)
	defer cancel()

	// Upload an object with storage.Writer.
	wc := s.bkt.Object(userID + ".gpg").NewWriter(ctx)
	wc.ChunkSize = 0

	if _, err := io.Copy(wc, buf); err != nil {
		return err
	}
	// Data can continue to be added to the file until the writer is closed.
	if err := wc.Close(); err != nil {
		return err
	}

	return nil
}

func (s *GoogleStorage) Download(userID string) ([]byte, error) {

	rc, err := s.bkt.Object(userID + ".gpg").NewReader(s.ctx)
	if err != nil {
		return []byte{}, err
	}
	defer rc.Close()
	slurp, err := ioutil.ReadAll(rc)
	if err != nil {
		return []byte{}, err
	}

	return slurp, nil

}
