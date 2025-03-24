package operator

import (
	"context"
	"fmt"
	"team00/client/calculator"
	"team00/client/model"
	"team00/client/saver"
	"team00/config"
	"team00/logger"
	pb "team00/server/device/server"
	"time"

	"google.golang.org/grpc"
)

type Operator interface {
	HandleValues(ctx context.Context, response grpc.ServerStreamingClient[pb.Message]) error
}

type operator struct {
	calculator calculator.Calculator
	saver      saver.Saver
	counter    int64
	threshold  int64
}

func NewOperator(cfg config.Operator) (Operator, error) {
	sv, err := saver.NewSaver(cfg.Saver)
	if err != nil {
		return nil, fmt.Errorf("failed to create new saver: %w", err)
	}

	calc := calculator.NewCalculator(cfg.Calculator)

	return &operator{
		calculator: calc,
		saver:      sv,
		threshold:  cfg.Threshold,
	}, nil
}

func (o *operator) HandleValues(ctx context.Context, response grpc.ServerStreamingClient[pb.Message]) error {
	msg, err := response.Recv()
	if err != nil {
		return fmt.Errorf("failed to recive message: %w", err)
	}

	t := time.Unix(msg.GetTimestamp(), 0)

	logger.GetLogger().Info().Msgf("Recived message. ID: %s, frequency: %f, timestamp: %s", msg.GetSessionId(), msg.GetFrequency(), t.String())

	calculated := o.calculator.Calculate(msg.GetFrequency())

	if err := o.saver.Save(ctx, model.SaverEntry{
		SessionId: msg.GetSessionId(),
		Frequency: msg.GetFrequency(),
		Timestamp: msg.GetTimestamp(),
	}); err != nil {
		return fmt.Errorf("failed to save: %w", err)
	}

	for _, v := range o.getSaverMessages(calculated) {
		if err := o.saver.Save(ctx, v); err != nil {
			return fmt.Errorf("failed to save: %w", err)
		}
	}

	o.counter++

	return nil
}

func (o *operator) getSaverMessages(calculated model.CalculatorRate) []model.SaverMessage {
	res := make([]model.SaverMessage, 0, 2)

	if o.counter%20 == 0 {
		res = append(res, model.SaverMean{
			Count:     o.counter,
			Mean:      calculated.Mean,
			Deviation: calculated.Deviation,
		})
	}

	if calculated.Anomaly.Exist {
		res = append(res, model.SaverAnomaly{
			Frequency: calculated.Anomaly.Frequency,
			Diff:      calculated.Anomaly.Diff,
		})
	}

	return res
}
