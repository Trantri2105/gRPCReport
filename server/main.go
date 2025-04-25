package main

import (
	"context"
	"fmt"
	"gRPC/proto/pb"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
)

type server struct {
	pb.UnimplementedMathServiceServer
}

func (*server) Sum(ctx context.Context, req *pb.SumRequest) (*pb.Response, error) {
	return &pb.Response{Num: req.FirstNum + req.SecondNum}, nil
}

func (*server) Avg(stream pb.MathService_AvgServer) error {
	cnt := int64(0)
	sum := int64(0)
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("Error receiving request: %v", err)
			return err
		}
		cnt += 1
		sum += req.Num
	}
	err := stream.SendAndClose(&pb.Response{Num: sum / cnt})
	if err != nil {
		log.Printf("Error sending response: %v", err)
		return err
	}
	return nil
}

func (*server) PrimeDivisor(req *pb.Request, res pb.MathService_PrimeDivisorServer) error {
	num := req.Num
	for num%2 == 0 {
		err := res.Send(&pb.Response{Num: 2})
		if err != nil {
			log.Printf("Error sending response: %v", err)
			return err
		}
		num /= 2
	}
	for i := int64(3); i*i <= num; i += 2 {
		if num%i == 0 {
			err := res.Send(&pb.Response{Num: i})
			if err != nil {
				log.Printf("Error sending response: %v", err)
				return err
			}
			num /= i
		}
	}
	if num != 1 {
		err := res.Send(&pb.Response{Num: num})
		if err != nil {
			log.Printf("Error sending response: %v", err)
		}
	}
	return nil
}

func (*server) Max(stream pb.MathService_MaxServer) error {
	mx := int64(0)
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Printf("Error receiving request: %v", err)
			return err
		}
		if mx < req.Num {
			mx = req.Num
		}
		err = stream.Send(&pb.Response{Num: mx})
		if err != nil {
			log.Printf("Error sending response: %v", err)
			return err
		}
	}
}

func NewServer() pb.MathServiceServer {
	return &server{}
}

func main() {
	port := 8080
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterMathServiceServer(grpcServer, NewServer())
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Printf("failed to start gRPC server: %v", err)
		return
	}
}
