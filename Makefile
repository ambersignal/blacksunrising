.PHONY: build run test lint clean

build:
	go build -o blacksunrising .

run: build
	./blacksunrising

test:
	go test ./...

lint:
	go tool golangci-lint run ./...

clean:
	rm -f blacksunrising
