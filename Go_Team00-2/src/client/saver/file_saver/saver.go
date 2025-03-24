package file_saver

import (
	"context"
	"errors"
	"fmt"
	"os"
	"team00/client/model"
	"team00/config"
	"team00/logger"
)

var ErrUnknownMessageType = errors.New("unknown message type")

type fileSaver struct {
	meanPath    string
	anomalyPath string
}

func NewFileSaver(cfg config.FileSaver) (*fileSaver, error) {
	if err := checkFile(cfg.MeanPath); err != nil {
		return nil, fmt.Errorf("failed to check file(%s): %w", cfg.MeanPath, err)
	}

	if err := checkFile(cfg.AnomalyPath); err != nil {
		return nil, fmt.Errorf("failed to check file(%s): %w", cfg.AnomalyPath, err)
	}

	return &fileSaver{
		meanPath:    cfg.MeanPath,
		anomalyPath: cfg.AnomalyPath,
	}, nil
}

func (s *fileSaver) Save(ctx context.Context, data model.SaverMessage) error {
	switch v := data.(type) {
	case model.SaverMean:
		if err := s.save(ctx, s.meanPath, v.String()); err != nil {
			return fmt.Errorf("failed to save mean: %w", err)
		}

	case model.SaverAnomaly:
		if err := s.save(ctx, s.anomalyPath, v.String()); err != nil {
			return fmt.Errorf("failed to save anomaly: %w", err)
		}

	default:
		return ErrUnknownMessageType
	}

	return nil
}

func (s *fileSaver) save(_ context.Context, path string, data string) error {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}

	defer func() {
		if err := f.Close(); err != nil {
			logger.GetLogger().Err(err).Msg("failed to close file")
		}
	}()

	if _, err := f.WriteString(data); err != nil {
		return fmt.Errorf("failed to write string: %w", err)
	}

	return nil
}

func checkFile(path string) error {
	if _, err := os.Stat(path); err == nil {
		return nil
	} else if errors.Is(err, os.ErrNotExist) {
		f, err := os.Create(path)
		if err != nil {
			return fmt.Errorf("failed to create file: %w", err)
		}

		if err := f.Close(); err != nil {
			return fmt.Errorf("failed to close file: %w", err)
		}

		logger.GetLogger().Debug().Msgf("file %s closed", path)

	} else {
		return fmt.Errorf("failed to get file stat: %w", err)
	}

	return nil
}

func (s *fileSaver) Close() error {
	return nil
}
