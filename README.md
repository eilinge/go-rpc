# go-rpc

    go grpc micro study

## protoc install

    https://www.jianshu.com/p/00be93ed230c

## protoc build

    protoc --proto_path=proto --proto_path=third_party --go_out=plugins=grpc:proto service.proto

## docker build

    docker build --tag=go-rpc .

## docker run

    docker run -it -p 8081:8081 go-rpc