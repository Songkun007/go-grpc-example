package main

import (
	"context"
	"log"

	pb "github.com/Songkun007/go-grpc-example/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const PORT = "9001"

func main() {
	// 支持 TLS 证书认证
	// credentials.NewServerTLSFromFile(certFile, serverNameOverride string)：
	// 根据客户端输入的证书文件和密钥构造 TLS 凭证。serverNameOverride 为服务名称
	c, err := credentials.NewClientTLSFromFile("../../conf/server.pem", "go-grpc-example")
	if err != nil {
		log.Fatalf("credentials.NewClientTLSFromFile err: %v", err)
	}

	conn, err := grpc.Dial(":"+PORT, grpc.WithTransportCredentials(c))
	if err != nil {
		log.Fatalf("grpc.Dial err: %v", err)
	}
	defer conn.Close()

	client := pb.NewSearchServiceClient(conn)
	resp, err := client.Search(context.Background(), &pb.SearchRequest{
		Request: "gRPC",
	})
	if err != nil {
		log.Fatalf("client.Search err: %v", err)
	}

	log.Printf("resp: %s", resp.GetResponse())
}
