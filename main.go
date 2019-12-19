package main

import (
	"context"
	"log"
	"net"
	"os"
	"sync"

	"go-rpc/proto"

	"google.golang.org/grpc"
	glog "google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/reflection"
)

var grpcLog glog.LoggerV2

func init() {
	grpcLog = glog.NewLoggerV2(os.Stdout, os.Stdout, os.Stdout)
}

type Connection struct {
	stream proto.Broadcast_CreateStreamServer // send message and server-side behavior of a streaming RPC.
	id     string
	name   string
	active bool
	error  chan error
}

type Server struct {
	Connection []*Connection
}

// create a chat zoom
func (s *Server) CreateStream(pconn *proto.Connect, stream proto.Broadcast_CreateStreamServer) error {
	conn := &Connection{
		stream: stream,
		id:     pconn.User.Id,
		name:   pconn.User.Name,
		active: true,
		error:  make(chan error),
	}

	s.Connection = append(s.Connection, conn)
	return <-conn.error
}

// for each connection that is active send message
func (s *Server) BroadcastMessage(ctx context.Context, msg *proto.Message) (*proto.Close, error) {
	wait := sync.WaitGroup{}
	done := make(chan int)

	// use goroutine for each connection that is active send message
	for _, conn := range s.Connection {
		wait.Add(1)

		go func(msg *proto.Message, conn *Connection) {
			defer wait.Done()

			if conn.active {
				err := conn.stream.Send(msg)
				grpcLog.Info("Sending message to: ", conn.name)

				// close connction if send error
				if err != nil {
					grpcLog.Errorf("Error with Stream %v - Error: %v", conn.name, err)
					conn.active = false
					conn.error <- err
				}
			}
		}(msg, conn)
	}

	go func() {
		wait.Wait()
		close(done)
	}()

	grpcLog.Infof("close service: %v", <-done)
	// proto.Close{} nothing(nil)
	return &proto.Close{}, nil
}

func main() {
	var connections []*Connection

	server := &Server{connections}

	grpcServer := grpc.NewServer() // 创建gRPC服务器

	listener, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatalf("error creating the server %v", err)
	}

	grpcLog.Info("Starting server at port: 8081")

	proto.RegisterBroadcastServer(grpcServer, server) // 在gRPC服务端注册服务
	reflection.Register(grpcServer)                   //在给定的gRPC服务器上注册服务器反射服务
	// Serve方法在lis上接受传入连接，为每个连接创建一个ServerTransport和server的goroutine。
	// 该goroutine读取gRPC请求，然后调用已注册的处理程序来响应它们。
	grpcServer.Serve(listener) // start server
}
