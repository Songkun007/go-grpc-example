package main

import (
	"context"
	"log"
	"net"

	pb "github.com/Songkun007/go-grpc-example/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type SearchService struct{}

func (s *SearchService) Search(ctx context.Context, r *pb.SearchRequest) (*pb.SearchResponse, error) {
	return &pb.SearchResponse{Response: r.GetRequest() + " Server"}, nil
}

const PORT = "9001"

func main() {
	// 支持 TLS 证书认证
	// credentials.NewServerTLSFromFile：根据服务端输入的证书文件和密钥构造 TLS 凭证
	c, err := credentials.NewServerTLSFromFile("../../conf/server.pem", "../../conf/server.key")
	if err != nil {
		log.Fatalf("credentials.NewServerTLSFromFile err: %v", err)
	}
	server := grpc.NewServer(grpc.Creds(c))
	// 将 SearchService（其包含需要被调用的服务端接口）注册到 gRPC Server 的内部注册中心。
	// 这样可以在接受到请求时，通过内部的服务发现，发现该服务端接口并转接进行逻辑处理
	pb.RegisterSearchServiceServer(server, &SearchService{})

	// 创建 Listen，监听 TCP 端口
	lis, err := net.Listen("tcp", ":"+PORT)
	if err != nil {
		log.Fatalf("net.Listen err: %v", err)
	}

	// gRPC Server 开始 lis.Accept，直到 Stop 或 GracefulStop
	server.Serve(lis)
}