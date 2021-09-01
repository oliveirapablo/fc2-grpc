package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/oliveirapablo/fc2-grpc/pb"
	"google.golang.org/grpc"
)

func main() {
	connection, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect to gRPC Server: %v", err)
	}
	defer connection.Close()

	client := pb.NewUserServiceClient(connection)

	// AddUser(client)
	// AddUserVerbose(client)
	// AddUsers(client)
	AddUserStreamBoth(client)
}

func AddUser(client pb.UserServiceClient) {
	req := &pb.User{
		Id:    "569",
		Name:  "Pablo",
		Email: "p@p.com",
	}

	res, err := client.AddUser(context.Background(), req)
	if err != nil {
		log.Fatalf("Could not make gRPC Request: %v", err)
	}

	fmt.Println(res)
}

func AddUserVerbose(client pb.UserServiceClient) {
	req := &pb.User{
		Id:    "569",
		Name:  "Pablo",
		Email: "p@p.com",
	}

	responseStream, err := client.AddUserVerbose(context.Background(), req)
	if err != nil {
		log.Fatalf("Could not make gRPC Request: %v", err)
	}

	for {
		stream, err := responseStream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("Could not receive the message stream: v%", err)
		}

		fmt.Println("Status: ", stream.Status, " - ", stream.GetUser())
	}

}

func AddUsers(client pb.UserServiceClient) {
	reqs := []*pb.User{
		&pb.User{
			Id:    "p1",
			Name:  "Pablo",
			Email: "pablocks@sant.com",
		},
		&pb.User{
			Id:    "p2",
			Name:  "Pablo2",
			Email: "pablocks2@sant.com",
		},
		&pb.User{
			Id:    "p3",
			Name:  "Pablo3",
			Email: "pablocks3@sant.com",
		},
		&pb.User{
			Id:    "p4",
			Name:  "Pablo4",
			Email: "pablocks4@sant.com",
		},
		&pb.User{
			Id:    "p5",
			Name:  "Pablo5",
			Email: "pablocks5@sant.com",
		},
	}
	stream, err := client.AddUsers(context.Background())
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	for _, req := range reqs {
		stream.Send(req)
		time.Sleep(time.Second * 3)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error receiving response: %v\n", err)
	}

	fmt.Println(res)
}

func AddUserStreamBoth(client pb.UserServiceClient) {
	stream, err := client.AddUserStreamBoth(context.Background())
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	reqs := []*pb.User{
		&pb.User{
			Id:    "p1",
			Name:  "Pablo",
			Email: "pablocks@sant.com",
		},
		&pb.User{
			Id:    "p2",
			Name:  "Pablo2",
			Email: "pablocks2@sant.com",
		},
		&pb.User{
			Id:    "p3",
			Name:  "Pablo3",
			Email: "pablocks3@sant.com",
		},
		&pb.User{
			Id:    "p4",
			Name:  "Pablo4",
			Email: "pablocks4@sant.com",
		},
		&pb.User{
			Id:    "p5",
			Name:  "Pablo5",
			Email: "pablocks5@sant.com",
		},
	}

	wait := make(chan int)

	go func() {
		for _, req := range reqs {
			fmt.Println("Sending user: ", req.Name)
			stream.Send(req)
			time.Sleep(time.Second * 2)
		}
		stream.CloseSend()
	}()

	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error receiving data: %v", err)
				break
			}
			fmt.Printf("Recebendo user %v com o status: %v", res.GetUser().GetName(), res.GetStatus())
		}
		close(wait)
	}()

	<-wait
}
