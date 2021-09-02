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