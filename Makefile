generate-swagger:
	swagger generate spec -w cmd/ -m -o assets/swagger/swagger.json

validate-swagger:
	swagger validate assets/swagger/swagger.json

build-app:
	go build -o bin/h24 ./cmd/

start-server:
	go run ./cmd/ server

run-server: generate-swagger validate-swagger start-server

run-aggregator:
	go run ./cmd/ aggregator

run-notificator:
	go run ./cmd/ notificator
