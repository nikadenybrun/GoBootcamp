package saver

import (
	"context"
	"errors"
	"fmt"
	"team00/client/model"
	"team00/client/saver/file_saver"
	"team00/client/saver/pg_saver"
	"team00/config"
)

var ErrUnknownSaverType = errors.New("unknown saver type")

const (
	FileSaver = "file_saver"
	DBSaver   = "database_saver"
)

type Saver interface {
	Save(ctx context.Context, data model.SaverMessage) error
	Close() error
}

func NewSaver(cfg config.Saver) (Saver, error) {
	switch cfg.Type {
	case FileSaver:
		s, err := file_saver.NewFileSaver(cfg.FileSaver)
		if err != nil {
			return nil, fmt.Errorf("failed to create new file saver: %w", err)
		}

		return s, nil

	case DBSaver:
		s, err := pg_saver.NewDatabaseSaver(cfg.DBSaver)
		if err != nil {
			return nil, fmt.Errorf("failed to create new database saver: %w", err)
		}

		return s, nil

	default:
		return nil, ErrUnknownSaverType
	}
}
