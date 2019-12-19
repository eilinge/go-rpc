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

## proto learn

    https://www.jianshu.com/p/2265f56805fa
    https://www.jianshu.com/p/774b38306c30
    https://www.jianshu.com/p/85e9cfa16247

    客户端 可以向 服务端 订阅 一个数据，服务端 就 可以利用 stream ，源源不断地 推送数据。

    ```
    message User {  # type User struct
        string id = 1; # GetId() string
        string name = 2; # GetName() string
    }

    service Broadcast { # type BroadcastClient interface
        rpc CreateStream(Connect) returns (stream Message); 
        # CreateStream(ctx context.Context, in *Connect, opts ...grpc.CallOption) (Broadcast_CreateStreamClient, error)
        rpc BroadcastMessage(Message) returns (Close); 
        # BroadcastMessage(ctx context.Context, in *Message, opts ...grpc.CallOption) (*Close, error)
    }
    ```