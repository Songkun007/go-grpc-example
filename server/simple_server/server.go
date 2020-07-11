package main

import (
	"context"
	"log"
	"net"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"runtime/debug"

	pb "github.com/Songkun007/go-grpc-example/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/codes"
	"github.com/grpc-ecosystem/go-grpc-middleware"
)

type SearchService struct{}

func (s *SearchService) Search(ctx context.Context, r *pb.SearchRequest) (*pb.SearchResponse, error) {
	return &pb.SearchResponse{Response: r.GetRequest() + " Server"}, nil
}

const PORT = "9001"

func main() {
	// 基于 CA 进行 TLS 认证
	// tls.LoadX509KeyPair(certFile, keyFile string)：从证书相关文件中读取和解析信息，得到证书公钥、密钥对
	cert, err := tls.LoadX509KeyPair("../../conf/server/server.pem", "../../conf/server/server.key")
	if err != nil {
		log.Fatalf("tls.LoadX509KeyPair err: %v", err)
	}

	certPool := x509.NewCertPool()
	ca, err := ioutil.ReadFile("../../conf/ca.pem")
	if err != nil {
		log.Fatalf("ioutil.ReadFile err: %v", err)
	}

	if ok := certPool.AppendCertsFromPEM(ca); !ok {
		log.Fatalf("certPool.AppendCertsFromPEM err")
	}

	// 在设置了 tls.RequireAndVerifyClientCert 模式的情况下，
	// Server 也会使用 CA 认证的根证书对 Client 端的证书进行可靠性、有效性等校验
	c := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    certPool,
	})

	opts := []grpc.ServerOption{
		grpc.Creds(c),
		grpc_middleware.WithUnaryServerChain(
			RecoveryInterceptor,
			LoggingInterceptor,
		),
	}


	server := grpc.NewServer(opts...)
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

// 拦截器

// logging
// 实现 RPC 方法的入参出参的日志输出
// ctx context.Context：请求上下文
// req interface{}：RPC 方法的请求参数
// info *UnaryServerInfo：RPC 方法的所有信息
// handler UnaryHandler：RPC 方法本身
func LoggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log.Printf("gRPC method: %s, %v", info.FullMethod, req)
	resp, err := handler(ctx, req)
	log.Printf("gRPC method: %s, %v", info.FullMethod, resp)
	return resp, err
}

// recover
// RPC 方法的异常保护和日志输出
func RecoveryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	defer func() {
		if e := recover(); e != nil {
			debug.PrintStack()
			err = status.Errorf(codes.Internal, "Panic err: %v", e)
		}
	}()

	return handler(ctx, req)
}