package storageproviders

type StorageProvider interface {
	Upload(userID string, pubkey []byte) error
	Download(userID string) ([]byte, error)
}
