FROM golang:latest

WORKDIR $GOPATH/src/github.com/gnixxyz/gin-sample
COPY . $GOPATH/src/github.com/gnixxyz/gin-sample

RUN go build .
