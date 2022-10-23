package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/picolloo/grpc-api-gateway/proto/protobuf"
)

func main() {
	conn, err := grpc.Dial(
		"localhost:9000",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		panic(err)
	}
	router := pb.NewRouterClient(conn)
	restSub, err := router.RestSubscribe(context.Background())
	rpcSub, err := router.RPCSubscribe(context.Background())

	client := http.DefaultServeMux

	client.HandleFunc("/rest", func(rw http.ResponseWriter, req *http.Request) {
		err := restSub.Send(&pb.RestRequest{
			Method: req.Method,
			Path:   req.URL.String(),
		})

		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, err.Error())))
			return
		}

		resp, err := restSub.Recv()
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, err.Error())))
			return
		}

		rw.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(rw).Encode(resp)
	})

	client.HandleFunc("/rpc", func(rw http.ResponseWriter, req *http.Request) {
		body, err := io.ReadAll(req.Body)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, err.Error())))
			return
		}

		var rpcReq pb.RPCRequest
		err = json.Unmarshal(body, &rpcReq)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, err.Error())))
			return
		}

		err = rpcSub.Send(&rpcReq)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, err.Error())))
			return
		}

		resp, err := rpcSub.Recv()
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, err.Error())))
			return
		}

		rw.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(rw).Encode(resp)
	})

	_ = http.ListenAndServe(":9001", client)
}
