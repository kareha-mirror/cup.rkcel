all: build

build:
	go build -o rkcel ./cmd/rkcel
	go build -o rkcel-calib ./cmd/rkcel-calib

clean:
	rm -f rkcel rkcel-calib

run:
	go run ./cmd/rkcel

fmt:
	go fmt ./...

test:
	go test ./...
