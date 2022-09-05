package config

type GoogleStorage struct {
	GoogleServiceAccount string `env:"GOOGLE_CREDENTIALS,required"`
	BucketName           string `env:"GOOGLE_STORAGE_BUCKET,required"`
}

type Server struct {
	PORT int `env:"PORT,default=50051"`
}
