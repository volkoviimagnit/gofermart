package db

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type ConnectionPostgres struct {
	ctx  context.Context
	dsn  string
	conn *pgxpool.Pool
}

func (s *ConnectionPostgres) Conn() *pgxpool.Pool {
	return s.conn
}

func NewConnectionPostgres(ctx context.Context, dsn string) *ConnectionPostgres {
	return &ConnectionPostgres{
		ctx: ctx,
		dsn: dsn,
	}
}

func (s *ConnectionPostgres) Exec(sql string, arguments ...any) error {
	_, errExecuting := s.conn.Exec(s.ctx, sql, arguments...)
	return errExecuting
}

func (s *ConnectionPostgres) Query(sql string, arguments ...any) pgx.Row {
	return s.conn.QueryRow(s.ctx, sql, arguments...)
}

// tryConnect
func (s *ConnectionPostgres) TryConnect() error {
	if s.conn != nil {
		return nil
	}

	pool, errorConnection := pgxpool.Connect(s.ctx, s.dsn)
	if errorConnection != nil {
		return errorConnection
	}
	s.conn = pool
	if s.conn == nil {
		return errors.New("не удалось установить соединение")
	}

	var d time.Time
	actualDate := time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)
	err := s.conn.QueryRow(context.Background(), "select $1::date", actualDate).Scan(&d)
	if err != nil {
		return errors.New("Unexpected failure on QueryRow Scan: " + err.Error())
	}
	if !actualDate.Equal(d) {
		return errors.New("Did not transcode date successfully: %v is not %v" + err.Error())
	}

	return nil
}
