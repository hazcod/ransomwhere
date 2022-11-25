FLAGS = -log=trace

all: run

test:
	go test -v ./...

clean:
	rm app || true

build:
	CGO_ENABLED=0 go build -o app -ldflags '-w -s -extldflags "-static"' ./cmd/...

run:
	go run ./cmd/... $(FLAGS)