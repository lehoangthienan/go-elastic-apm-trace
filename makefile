setup:
	docker-compose up

init:
	cd svc-a && cat .env.example > .env
	cd svc-b && cat .env.example > .env
	go mod tidy

dev-svc-a:
	cd svc-a && source .env && go run main.go

dev-svc-b:
	cd svc-b && source .env && go run main.go

gen-proto:
	go get -u google.golang.org/protobuf/cmd/protoc-gen-go && go install google.golang.org/protobuf/cmd/protoc-gen-go
	go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc && go install google.golang.org/grpc/cmd/protoc-gen-go-grpc
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/user.proto
