all: build

build:
	go build -o rkcel ./cmd/rkcel

clean:
	rm -f rkcel

run:
	go run ./cmd/rkcel

fmt:
	go fmt ./...

test:
	go test ./...
