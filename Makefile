all: load fmt test build

load:
	go get ./...

fmt:
	go fmt ./...

test:
	go test ./...

build:
	GOOS=linux go build

clean:
	rm serverless-blog-back

run:
	sam local start-api