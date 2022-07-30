package pkg

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"time"

	storage "cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
)

type Storage struct {
	ctx    context.Context
	client *storage.Client
	bkt    *storage.BucketHandle
}

func NewStorageClient() *Storage {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("Unable to initialize the storage client %v ", err)
	}

	bkt := client.Bucket("oss-pubkeys")

	return &Storage{
		ctx:    ctx,
		client: client,
		bkt:    bkt,
	}
}

func (s *Storage) Close() {
	s.client.Close()
}

func (s *Storage) Upload(userID string, pubkey []byte) error {
	buf := bytes.NewBuffer(pubkey)

	ctx, cancel := context.WithTimeout(s.ctx, time.Second*50)
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

func (s *Storage) Download(userID string) ([]byte, error) {

	ctx, cancel := context.WithTimeout(s.ctx, time.Second*50)

	query := &storage.Query{Prefix: userID}

	var names []string
	it := s.bkt.Objects(ctx, query)
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		names = append(names, attrs.Name)
	}
	defer cancel()

	fmt.Println("objects : ", names)
	return []byte{}, nil
}
