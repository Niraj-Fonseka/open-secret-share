


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

build-go:
	@echo "----- building go binary -----"
	@go build -o server .

build-docker:
	@echo "----- building docker image -----"
	@docker build \
	--build-arg google_creds=$(google_creds) \
	--build-arg auth_key=$(auth_key) \
	--build-arg google_storage_bucket=$(google_storage_bucket) \
	-t oss-key-server:latest .