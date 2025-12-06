BIN := bin/daytime

build: bin
	go build -o $(BIN) cmd/main.go

run: build
	$(BIN) $(ARGS)

test:
	go test ./...

bin:
	mkdir -p bin

tidy:
	go mod tidy
