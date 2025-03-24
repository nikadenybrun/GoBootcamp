package pg_saver

import (
	"context"
	"errors"
	"fmt"
	"team00/client/model"
	"team00/config"
	"team00/logger"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

var ErrUnknownMessageType = errors.New("unknown message type")

type pgSaver struct {
	db *pg.DB
}

func (s *pgSaver) migrateSchema() error {
	schemas := []interface{}{
		(*model.SaverMean)(nil),
		(*model.SaverAnomaly)(nil),
		(*model.SaverEntry)(nil),
	}

	op := orm.CreateTableOptions{IfNotExists: true}

	for _, schema := range schemas {
		if err := s.db.Model(schema).CreateTable(&op); err != nil {
			return fmt.Errorf("failed to create table: %w", err)
		}
	}

	return nil
}

func NewDatabaseSaver(cfg config.DBSaver) (*pgSaver, error) {
	conn := pg.Connect(&pg.Options{
		User:     cfg.User,
		Password: cfg.Password,
		Database: cfg.Database,
	})

	if err := conn.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	obj := &pgSaver{
		db: conn,
	}

	if err := obj.migrateSchema(); err != nil {
		return nil, fmt.Errorf("failed to migrate schema: %w", err)
	}

	return obj, nil
}

func (s *pgSaver) Save(_ context.Context, msg model.SaverMessage) error {
	if v, ok := msg.(model.SaverAnomaly); ok {
		if err := saveMessage(s.db, &v); err != nil {
			return fmt.Errorf("failed to save message: %w", err)
		}

		return nil
	}

	if v, ok := msg.(model.SaverMean); ok {
		if err := saveMessage(s.db, &v); err != nil {
			return fmt.Errorf("failed to save message: %w", err)
		}

		return nil
	}

	if v, ok := msg.(model.SaverEntry); ok {
		if err := saveMessage(s.db, &v); err != nil {
			return fmt.Errorf("failed to save message: %w", err)
		}

		return nil
	}

	logger.GetLogger().Debug().Msg(msg.String())

	return ErrUnknownMessageType
}

func saveMessage[T any](db *pg.DB, msg T) error {
	if _, err := db.Model(msg).Insert(); err != nil {
		return fmt.Errorf("failed to insert: %w", err)
	}

	return nil
}

func (s *pgSaver) Close() error {
	if err := s.db.Close(); err != nil {
		return fmt.Errorf("failed to close db: %w", err)
	}

	return nil
}
