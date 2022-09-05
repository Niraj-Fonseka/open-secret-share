


ifndef AUTH_KEY
$(error AUTH_KEY is not set)
endif

ifndef GOOGLE_CREDENTIALS
$(error GOOGLE_CREDENTIALS is not set)
endif

ifndef GOOGLE_STORAGE_BUCKET
$(error GOOGLE_STORAGE_BUCKET is not set)
endif


server_auth_key  = ${AUTH_KEY}
server_google_cred = ${GOOGLE_CREDENTIALS}
server_bucket_name  = ${GOOGLE_STORAGE_BUCKET}

build-go:
	@echo "----- building go binary -----"
	@go build -o server .

build-docker:
	@echo "----- building docker image -----"
	@docker build \
	--build-arg google_creds=$(server_google_cred) \
	--build-arg auth_key=$(server_auth_key) \
	--build-arg google_storage_bucket=$(server_bucket_name) \
	-t oss-key-server:latest .