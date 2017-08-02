package main

import (
	"net"
	pb "dirtyfilter/proto"


	log "github.com/Sirupsen/logrus"
	"google.golang.org/grpc"
)


func main() {
	lis, err := net.Listen("tcp", ":50002")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	// 注册服务
	s := grpc.NewServer()
	ins := &server{}
	ins.init()
	// 开始服务
	pb.RegisterWordFilterServiceServer(s, ins)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
