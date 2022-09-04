


build-go:
	go build -o server .

build-docker:
	
	docker build -t oss-key-server:latest --build-arg google_creds=$(GOOGLE_CREDENTIALS) --build-arg auth_key=$(AUTH_KEY) --build-arg google_storage_bucket=$(GOOGLE_STORAGE_BUCKET) .