FROM golang:alpine as build-env

ENV GO111MODULE=on

RUN apk update && apk add bash ca-certificates git gcc g++ libc-dev

RUN mkdir /go-rpc
RUN mkdir -p /go-rpc/proto

WORKDIR /go-prc

COPY ./vendor /go-rpc/vendor
COPY ./proto/service.pb.go /go-rpc/proto
COPY ./main.go /go-rpc

COPY go.mod .
COPY go.sum .

RUN go mod download

RUN go build -o go-rpc .

CMD ./go-rpc