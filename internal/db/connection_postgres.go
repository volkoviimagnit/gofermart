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
	errConnecting := s.TryConnect()
	if errConnecting != nil {
		return errConnecting
	}
	_, errExecuting := s.conn.Exec(s.ctx, sql, arguments...)
	return errExecuting
}

func (s *ConnectionPostgres) QueryRow(sql string, arguments ...any) (pgx.Row, error) {
	return s.conn.QueryRow(s.ctx, sql, arguments...), nil
}

func (s *ConnectionPostgres) Query(sql string, arguments ...any) (pgx.Rows, error) {
	return s.conn.Query(s.ctx, sql, arguments...)
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
		return errors.New("did not transcode date successfully")
	}

	return nil
}
