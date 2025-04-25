package main

import (
	"context"
	"gRPC/proto/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"log"
	"sync"
	"time"
)

func main() {
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	conn, err := grpc.NewClient("localhost:8080", opts...)
	if err != nil {
		log.Fatalf("failed to create gRPC connection: %v", err)
	}
	defer conn.Close()
	client := pb.NewMathServiceClient(conn)
	Max(client, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
}

func Max(client pb.MathServiceClient, ar ...int64) {
	stream, err := client.Max(context.Background())
	if err != nil {
		log.Printf("failed to invoke max rpc: %v", err)
		return
	}
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for _, num := range ar {
			err := stream.Send(&pb.Request{Num: num})
			if err != nil {
				log.Printf("failed to send request: %v", err)
				return
			}
		}
		err := stream.CloseSend()
		if err != nil {
			log.Printf("failed to close send stream: %v", err)
			return
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				return
			}
			if err != nil {
				log.Printf("failed to receive response: %v", err)
				return
			}
			log.Println(res.Num)
		}
	}()
	wg.Wait()
}

func PrimeDivisor(client pb.MathServiceClient, num int64) {
	stream, err := client.PrimeDivisor(context.Background(), &pb.Request{Num: num})
	if err != nil {
		log.Printf("failed to send request: %v", err)
		return
	}
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			return
		}
		if err != nil {
			log.Printf("failed to receive response: %v", err)
			return
		}
		log.Println(res.Num)
	}
}

func Avg(client pb.MathServiceClient, ar ...int64) {
	stream, err := client.Avg(context.Background())
	if err != nil {
		log.Printf("failed to invoke avg request: %v", err)
		return
	}
	for _, e := range ar {
		err = stream.Send(&pb.Request{Num: e})
		if err != nil {
			log.Printf("failed to send request: %v", err)
			return
		}
	}
	resp, err := stream.CloseAndRecv()
	if err != nil {
		log.Printf("failed to receive response: %v", err)
		return
	}
	log.Println(resp.Num)
}

func Sum(firstNumber, secondNumber int64, client pb.MathServiceClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	res, err := client.Sum(ctx, &pb.SumRequest{FirstNum: firstNumber, SecondNum: secondNumber})
	if err != nil {
		log.Printf("failed to send sum request: %v", err)
		return
	}
	log.Println(res.Num)
}
