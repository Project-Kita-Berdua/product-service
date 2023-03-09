proto: 
	protoc pkg/pb/*.proto --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=require_unimplemented_servers=false
	
fix-proto:
	export GO_PATH=~/go
	export PATH=$PATH:/$GO_PATH/bin

server: 
	go run cmd/main.go