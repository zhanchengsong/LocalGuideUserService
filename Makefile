build:
	swag init -g main.go
	go build

run:
	swag init -g main.go
	go run main.go