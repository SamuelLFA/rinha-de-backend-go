BINARY_NAME=app
DEBUG_FLAGS=-gcflags "all=-N -l"
GOFLAGS=-v
MAIN_FILE=cmd/main.go
ENV_VARS=DB_HOST=localhost DB_PORT=5432 DB_USER=postgres DB_PASSWORD=postgres DB_NAME=postgres DB_SSLMODE=disable

build:
	go build -o bin/$(BINARY_NAME) $(GOFLAGS) $(MAIN_FILE)
test:
	go test ./internal/...
clean:
	go clean
	rm -f bin/$(BINARY_NAME)
deps:
	go mod tidy
	go mod vendor
debug:
	go build $(DEBUG_FLAGS) -o bin/$(BINARY_NAME) $(GOFLAGS) $(MAIN_FILE)
run:
	$(ENV_VARS) go run $(MAIN_FILE)
dev-up:
	docker-compose -f docker/docker-compose.dev.yaml up -d