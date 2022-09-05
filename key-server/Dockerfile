#first stage - builder
FROM golang:1.18.0-stretch as builder

COPY . /oss-key-server
WORKDIR /oss-key-server

ENV GO111MODULE=on


RUN CGO_ENABLED=0 GOOS=linux go build -o server 


#second stage 
FROM alpine:latest

ARG google_credentials
ARG auth_key
ARG google_storage_bucket

ENV GOOGLE_CREDENTIALS=${google_creds}
ENV AUTH_KEY=${auth_key}
ENV GOOGLE_STORAGE_BUCKET=${google_storage_bucket}

RUN apk add --no-cache tzdata

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

WORKDIR /root/

COPY --from=builder /oss-key-server .

CMD ["./server"]