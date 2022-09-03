package storageproviders

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"time"

	storage "cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
)

type GoogleStorage struct {
	ctx    context.Context
	client *storage.Client
	bkt    *storage.BucketHandle
}

func NewGoogleStorageClient() *GoogleStorage {
	ctx := context.Background()

	//client, err := storage.NewClient(ctx, option.WithCredentialsJSON([]byte(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"))))
	client, err := storage.NewClient(ctx)

	if err != nil {
		log.Fatalf("Unable to initialize the storage client %v ", err)
	}

	bkt := client.Bucket("oss-pubkeys")

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
	wc.ChunkSize = 0 // note retries are not supported for chunk size 0.

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

	ctx, cancel := context.WithTimeout(s.ctx, time.Second*120)

	query := &storage.Query{Prefix: userID}

	var names []string
	it := s.bkt.Objects(ctx, query)
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return []byte{}, err
		}
		names = append(names, attrs.Name)
	}
	defer cancel()

	if len(names) == 0 {
		return []byte{}, fmt.Errorf("reciever not found ")
	}
	rc, err := s.bkt.Object(names[0]).NewReader(s.ctx)
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
