package main

import (
	"context"
	"io"
	"log"

	pb "github.com/Songkun007/go-grpc-example/proto"
	"google.golang.org/grpc"
)

const (
	PORT = "9002"
)

func main() {
	// 建立连接
	conn, err := grpc.Dial(":" + PORT, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("grpc.Dial err: %v", err)
	}

	defer conn.Close()

	client := pb.NewStreamServiceClient(conn)

	err = printLists(client, &pb.StreamRequest{Pt: &pb.StreamPoint{Name: "gRPC Stream Client: List", Value: 2018}})
	if err != nil {
		log.Fatalf("printLists.err: %v", err)
	}

	err = printRecord(client, &pb.StreamRequest{Pt: &pb.StreamPoint{Name: "gRPC Stream Client: Record", Value: 2018}})
	if err != nil {
		log.Fatalf("printRecord.err: %v", err)
	}

	err = printRoute(client, &pb.StreamRequest{Pt: &pb.StreamPoint{Name: "gRPC Stream Client: Route", Value: 2018}})
	if err != nil {
		log.Fatalf("printRoute.err: %v", err)
	}
}

func printLists(client pb.StreamServiceClient, r *pb.StreamRequest) error {
	stream, err := client.List(context.Background(), r)
	if err != nil {
		return err
	}

	for {
		resp, err := stream.Recv()
		// io.EOF 是在没有任何可读取的内容时触发，
		// 比如某文件Reader对象，文件本身为空，或者读取若干次后，文件指针指向了末尾，调用Read都会触发EOF。
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		log.Printf("resp: pj.name: %s, pt.value: %d", resp.Pt.Name, resp.Pt.Value)
	}

	return nil
}

func printRecord(client pb.StreamServiceClient, r *pb.StreamRequest) error {
	stream, err := client.Record(context.Background())
	if err != nil {
		return err
	}

	for n := 0; n < 6; n++ {
		err := stream.Send(r)
		if err != nil {
			return err
		}
	}

	resp, err := stream.CloseAndRecv()
	if err != nil {
		return err
	}

	log.Printf("resp: pj.name: %s, pt.value: %d", resp.Pt.Name, resp.Pt.Value)

	return nil
}

func printRoute(client pb.StreamServiceClient, r *pb.StreamRequest) error {
	stream, err := client.Route(context.Background())
	if err != nil {
		return err
	}

	for n := 0; n <= 6; n++ {
		err = stream.Send(r)
		if err != nil {
			return err
		}

		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		log.Printf("resp: pj.name: %s, pt.value: %d", resp.Pt.Name, resp.Pt.Value)
	}

	return nil
}