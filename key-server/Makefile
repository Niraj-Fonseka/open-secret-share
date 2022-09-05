


ifndef AUTH_KEY
$(error AUTH_KEY is not set)
endif

ifndef GOOGLE_CREDENTIALS
$(error GOOGLE_CREDENTIALS is not set)
endif

ifndef GOOGLE_STORAGE_BUCKET
$(error GOOGLE_STORAGE_BUCKET is not set)
endif

auth_key=${AUTH_KEY}
google_creds=${GOOGLE_CREDENTIALS}
google_storage_bucket=${GOOGLE_STORAGE_BUCKET}
port=${PORT}

ifndef PORT
port=50051
endif


go: build-go run-go

docker: build-docker run-docker

run-docker:
	@docker run -p $(port):$(port) oss-key-server:latest

run-go:
	./oss-key-server

build-go:
	@echo "----- building go binary -----"
	@go build -o oss-key-server .

build-docker:
	@echo "----- building docker image -----"
	@docker build \
	--build-arg google_creds=$(google_creds) \
	--build-arg auth_key=$(auth_key) \
	--build-arg google_storage_bucket=$(google_storage_bucket) \
	--build-arg port=$(port) \
	-t oss-key-server:latest .