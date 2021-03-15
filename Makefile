BINARY=./bin/parser
start: build
	@$(BINARY)

build:
	go build -o ./bin/parser