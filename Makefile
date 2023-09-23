run:
	go run ./main.go

build:
	GOOS=linux GOARCH=amd64 go build -o ./dist/cardstore

migrate-up:
	go run ./main.go --run=migrate-up

migrate-down-all:
	go run ./main.go --run=migrate-down-all

delete-expired-sessions:
	go run ./main.go --run=delete-expired-sessions

add-keys:
	go run ./main.go --run=add-keys

encrypt-unencrypted:
	go run ./main.go --run=encrypt-unencrypted
