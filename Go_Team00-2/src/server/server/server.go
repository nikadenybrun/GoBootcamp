package main

import (
	"net"
	"os"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"team00/config"
	"team00/logger"
	pb "team00/server/device/server"
	"team00/server/server/generator"
)

type deviceServer struct {
	generator generator.Generator
	pb.UnimplementedDeviceServer
}

func newDeviceServer() *deviceServer {
	return &deviceServer{
		generator: generator.NewGenerator(),
	}
}

func (s *deviceServer) Connect(req *pb.Empty, stream grpc.ServerStreamingServer[pb.Message]) error {
	id, err := uuid.NewUUID()
	if err != nil {
		return status.Error(codes.Internal, "unable to generate uuid")
	}

	msg := pb.Message{
		SessionId: id.String(),
		Frequency: s.generator.Generate(),
		Timestamp: time.Now().UTC().Unix(),
	}

	if err := stream.Send(&msg); err != nil {
		return status.Error(codes.Internal, "unable to send")
	}

	return nil
}

func main() {
	logger.InitLogger("ALIEN TROOPER", config.Logger{
		Level: "info",
	})

	cfg := config.Client{
		Host: "127.0.01",
		Port: "50051",
	}

	lis, err := net.Listen("tcp", net.JoinHostPort(cfg.Host, cfg.Port))
	if err != nil {
		logger.GetLogger().Err(err).Msg("failed to listen port")

		os.Exit(1)
	}

	srv := grpc.NewServer()
	pb.RegisterDeviceServer(srv, newDeviceServer())

	logger.GetLogger().Info().Msgf("gRPC server listening on %s", lis.Addr())
	srv.Serve(lis)
}
