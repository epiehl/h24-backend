generate-swagger:
	swagger generate spec -w cmd/server -m -o assets/swagger/swagger.json

validate-swagger:
	swagger validate assets/swagger/swagger.json

build-server:
	go build -o bin/server ./cmd/server

start-server:
	go run ./cmd/server

run-server: generate-swagger validate-swagger start-server

run-aggregator:
	go run ./cmd/aggregator

build-aggregator:
	go build -o bin/aggregator ./cmd/aggregator

run-notificator:
	go run ./cmd/notificator

build-aggregator:
	go build -o bin/notificator ./cmd/notificator
