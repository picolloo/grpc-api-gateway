package main

import (
	"net"
	"net/http"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "github.com/picolloo/grpc-api-gateway/proto/protobuf"
)

type Router struct {
	pb.UnimplementedRouterServer
}

func NewRouter() *Router {
	return &Router{}
}

func (r *Router) RestSubscribe(srv pb.Router_RestSubscribeServer) error {
	for {
		req, err := srv.Recv()
		if err != nil {
			return err
		}

		err = srv.Send(&pb.RestResponse{
			Message: req.Method,
			Status:  http.StatusOK,
		})

		if err != nil {
			return err
		}
	}
}

func (r *Router) RPCSubscribe(srv pb.Router_RPCSubscribeServer) error {
	for {
		req, err := srv.Recv()
		if err != nil {
			return err
		}

		err = srv.Send(&pb.RPCResponse{
			Id:      req.Id,
			Jsonrpc: req.Jsonrpc,
			Result:  req.Params,
		})

		if err != nil {
			return err
		}
	}
}

func main() {
	server := grpc.NewServer()
	router := NewRouter()
	pb.RegisterRouterServer(server, router)

	l, err := net.Listen("tcp", ":9000")
	if err != nil {
		panic(err)
	}

	reflection.Register(server)

	err = server.Serve(l)
	if err != nil {
		panic(err)
	}
}
