

all: generate cp

generate:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative oss.proto

cp: 
	cp oss.pb.go ./../key-server/protobuf/oss.pb.go
	cp oss_grpc.pb.go ./../key-server/protobuf/oss_grpc.pb.go
	cp oss.pb.go ./../oss/protobuf/oss.pb.go
	cp oss_grpc.pb.go ./../oss/protobuf/oss_grpc.pb.go