generate-wire:
	wire ./internal

generate-swagger:
	swagger generate spec -w cmd/ -m -o assets/swagger/swagger.json

validate-swagger:
	swagger validate assets/swagger/swagger.json

build-app:
	go build -o bin/h24 ./cmd/

start-server:
	go run ./cmd/ server

run-server: generate-wire generate-swagger validate-swagger start-server
run-aggregator: generate-wire start-aggregator
run-notificator: generate-wire start-notificator

start-aggregator:
	go run ./cmd/ aggregator

start-notificator:
	go run ./cmd/ notificator
