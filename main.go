package main

import (
	"context"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	pb "github.com/myuser/myrepo/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
)

type server struct {
	pb.UnimplementedProductServiceServer
}

func (s *server) DeleteFoo(ctx context.Context, req *pb.DeleteFooRequest) (*emptypb.Empty, error) {
	// Print the received soft_delete value
	println("soft_delete:", req.SoftDelete)
	println("soft_delete:", *req.SoftDelete)
	println("soft_delete:", &req.SoftDelete)
	println("================================= \n")

	return &emptypb.Empty{}, nil
}

func main() {
	grpcServer := grpc.NewServer()
	pb.RegisterProductServiceServer(grpcServer, &server{})

	// Start gRPC server
	go func() {
		listen, err := net.Listen("tcp", ":9090")
		if err != nil {
			panic(err)
		}
		reflection.Register(grpcServer)
		if err := grpcServer.Serve(listen); err != nil {
			panic(err)
		}
	}()

	// Start HTTP Gateway server
	mux := runtime.NewServeMux()
	err := pb.RegisterProductServiceHandlerFromEndpoint(context.Background(), mux, ":9090", []grpc.DialOption{grpc.WithInsecure()})
	if err != nil {
		panic(err)
	}

	if err := http.ListenAndServe(":8080", mux); err != nil {
		panic(err)
	}
}
