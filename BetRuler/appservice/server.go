package main

import (
	// "github.com/golang/protobuf/proto"
	"net"

	"google.golang.org/grpc"

	"BetRuler/appservice/api"
	"BetRuler/pb"
)

func main() {
	server := grpc.NewServer()

	pb.RegisterLotteryQueryServer(server, api.NewLotteryQueryService())

	listen, err := net.Listen("tcp", ":8090")
	if err != nil {
		panic(err.Error())
	}

	server.Serve(listen)
}
