package main

import (
	pb "dirtyfilter/proto"
	"fmt"
	"net"

	"github.com/vrischmann/envconfig"
	"google.golang.org/grpc"
	"os"
	"log"
)

func main() {
	log.SetOutput(os.Stdout)
	var config Config
	if err := envconfig.Init(&config); err != nil {
		panic(fmt.Sprintf("parse env failed: %s", err))
	}
	lis, err := net.Listen("tcp", ":50002")
	if err != nil {
		panic(fmt.Sprintf("failed to listen: %v", err))
	}
	log.Printf("dirtyfile host: %s", config.Host)
	log.Printf("dirtyfile port: %d", config.Port)
	// 注册服务
	s := grpc.NewServer()
	ins := &server{}
	ins.init()
	// 开始服务
	pb.RegisterWordFilterServiceServer(s, ins)
	log.Printf("dirtyfile server start success!")
	if err := s.Serve(lis); err != nil {
		log.Printf("555555555")
		panic(fmt.Sprintf("failed to serve: %v", err))
	}
}
