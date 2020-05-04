package main

import (
	"log"
	"net"

	pb "github.com/Songkun007/go-grpc-example/proto"
	"google.golang.org/grpc"
)

type StreamService struct {}

const (
	PORT = "9002"
)

func main() {
	// 创建一个gRPG对象
	server := grpc.NewServer()

	// 注册rpc服务
	pb.RegisterStreamServiceServer(server, &StreamService{})

	// 监听 TCP 端口
	lis, err := net.Listen("tcp", ":"+PORT)
	if err != nil {
		log.Fatalf("net.Listen err: %v", err)
	}

	server.Serve(lis)
}

func (s *StreamService) List(r *pb.StreamRequest, stream pb.StreamService_ListServer) error {
	for n := 0; n <= 6; n++ {
		err := stream.Send(&pb.StreamResponse{
			Pt: &pb.StreamPoint{
				Name:  r.Pt.Name,
				Value: r.Pt.Value + int32(n),
			},
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *StreamService) Record(stream pb.StreamService_RecordServer) error {
	return nil
}

func (s *StreamService) Route(stream pb.StreamService_RouteServer) error {
	return nil
}



