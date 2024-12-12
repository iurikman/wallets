build:
	@echo 'Building the project...'
	go build -o wallets/cmd/service/main

run: build
	@echo 'Running the project...'
	./wallets/cmd/service/main

lint:
	@echo 'Linting the project...'
	gofumpt -w .
	go mod tidy
	golangci-lint run --fix
	golangci-lint run --config .golangci.yaml

test: up
	go test -v ./...

up:
	docker compose up -d

down:
	docker compose down
