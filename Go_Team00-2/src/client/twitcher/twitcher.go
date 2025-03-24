package twitcher

import (
	"context"
	"errors"
	"fmt"
	"team00/client/operator"
	"team00/logger"
	pb "team00/server/device/server"

	"google.golang.org/grpc"
)

var ErrEmptyHandler = errors.New("empty handler")

type Twitcher interface {
	Twitch(ctx context.Context) error
	Close() error
}

type HandlerType func(response grpc.ServerStreamingClient[pb.Message]) error

type twitcher struct {
	conn     *grpc.ClientConn
	operator operator.Operator
	client   pb.DeviceClient
}

func NewTwitcher(conn *grpc.ClientConn, op operator.Operator) (Twitcher, error) {
	cl := pb.NewDeviceClient(conn)

	return &twitcher{
		client:   cl,
		operator: op,
	}, nil
}

func (t *twitcher) Twitch(ctx context.Context) error {
	logger.GetLogger().Info().Msg(`Press any key for next value`)
	logger.GetLogger().Info().Msg(`Press "Q" for exit`)
	for {
		if isBreak() {
			break
		}

		response, err := t.client.Connect(ctx, &pb.Empty{})
		if err != nil {
			return fmt.Errorf("failed to connect: %w", err)
		}

		if err := t.operator.HandleValues(ctx, response); err != nil {
			return fmt.Errorf("failed to handle: %w", err)
		}
	}

	return nil
}

func (t *twitcher) Close() error {
	if err := t.conn.Close(); err != nil {
		return fmt.Errorf("failed to close grpc connection: %w", err)
	}

	return nil
}

func isBreak() bool {
	var key byte

	fmt.Scanf("%c", &key)

	if key == 'Q' || key == 'q' {
		return true
	}

	return false
}
