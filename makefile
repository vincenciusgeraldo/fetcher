dep:
	GO111MODULE=on go mod download
	GO111MODULE=on go mod vendor

build:
	GO111MODULE=on go build -o bin/fetcher main.go

test:
	GO111MODULE=on go test ./...