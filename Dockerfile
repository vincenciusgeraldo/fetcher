FROM golang:1.16

WORKDIR /go/src/app
COPY . .

RUN go mod download && \
	go mod vendor
RUN go build -o bin/fetcher main.go

ENTRYPOINT ["bin/fetcher"]