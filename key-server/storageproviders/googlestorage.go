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

	//fetching encoded GOOGLE_APPLICATION_CREDENTIALS
	// decoded, err := base64.RawStdEncoding.DecodeString(googleConfig.GoogleServiceAccount)
	// if err != nil {
	// 	log.Fatalf("Unable to decode google service account %s ", err.Error())
	// }

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

// func (s *GoogleStorage) Peek(filename string) (bool, error) {
// 	ctx, cancel := context.WithTimeout(s.ctx, time.Second*120)

// 	defer cancel()

// 	bucket := s.bkt.
// }

func (s *GoogleStorage) Download(userID string) ([]byte, error) {

	//ctx, cancel := context.WithTimeout(s.ctx, time.Second*120)

	//query := &storage.Query{Prefix: userID}

	// var names []string
	// it := s.bkt.Object(ctx, query)
	// for {
	// 	attrs, err := it.Next()
	// 	if err == iterator.Done {
	// 		break
	// 	}
	// 	if err != nil {
	// 		return []byte{}, err
	// 	}
	// 	names = append(names, attrs.Name)
	// }
	// defer cancel()

	// if len(names) == 0 {
	// 	return []byte{}, fmt.Errorf("reciever not found ")
	// }
	rc, err := s.bkt.Object(userID).NewReader(s.ctx)
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
