package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"google.golang.org/grpc"

	pb "team00/server/device/server" // Import the generated code
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
	}
	defer conn.Close()

	client := pb.NewDeviceClient(conn)

	ctx := context.Background()
	stream, err := client.Connect(ctx, &pb.Empty{})
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	file, err := os.Create("log.txt")
	if err != nil {
		log.Println("did not manage to open the file")
	}
	defer file.Close()
	for i := 0; i < 101; i++ {
		msg, err := stream.Recv()
		if err != nil {
			log.Printf("error receiving message: %v", err)
			break
		}

		s := fmt.Sprintf("session_id: %s frequency: %f timestamp: %d", msg.SessionId, msg.Frequency, msg.Timestamp)
		_, err = file.WriteString(s + "\n")

		if err != nil {
			log.Println(err)
		}
	}
}
