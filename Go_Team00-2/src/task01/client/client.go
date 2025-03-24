package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math"
	"os"

	"google.golang.org/grpc"

	pb "team00/server/device/server" // Import the generated code
)

type stats struct {
	mean     float64
	variance float64
	count    int
}

func newStats() *stats {
	return &stats{}
}

func (s *stats) lesgo(freq []float64, k float64) {
	file, err := os.Create("journal.txt")
	if err != nil {
		log.Println(err)
	}
	defer file.Close()

	for idx, value := range freq {
		s.count++
		delta := value - s.mean
		s.mean += delta / float64(s.count)
		s.variance += delta * (value - s.mean)

		if idx%20 == 0 {
			s.writeJournal(file, idx)
		}

		if idx > 20 {
			s.checkAnomaly(file, k, value, idx)
		}
	}
}

func (s *stats) writeJournal(file *os.File, idx int) {
	str := fmt.Sprintf("idx: %d mean: %f, STD: %f", idx, s.mean, s.stddev())
	_, err := file.WriteString(str + "\n")
	if err != nil {
		log.Println(err)
	}
}

func (s *stats) checkAnomaly(file *os.File, k float64, num float64, idx int) {
	anomalyFlag := k * s.stddev()
	num = math.Abs(num)
	value := num - s.mean
	if value > anomalyFlag || (-1*value) < (-1*anomalyFlag) {
		str := fmt.Sprintf("idx: %d anomaly: %f", idx, num)
		_, err := file.WriteString(str + "\n")
		if err != nil {
			log.Println(err)
		}
	}
}

func (s *stats) stddev() float64 {
	if s.count < 2 {
		return 0
	}

	return math.Sqrt(s.variance / float64(s.count-1))
}

func main() {
	k := flag.Float64("k", 1, "anomaly")
	flag.Parse()
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

	var freq []float64
	for i := 0; i < 101; i++ {
		msg, err := stream.Recv()
		if err != nil {
			log.Printf("error receiving message: %v", err)
		}

		freq = append(freq, msg.Frequency)
	}

	s := newStats()
	s.lesgo(freq, *k)
}
