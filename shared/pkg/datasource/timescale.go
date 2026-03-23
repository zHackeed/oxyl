package datasource

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"zhacked.me/oxyl/shared/pkg/variables"
)

type TimescaleConnection struct {
	conn *pgxpool.Pool
}

func NewTimescaleConnection(ctx context.Context) (*TimescaleConnection, error) {
	connectionCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	hostValues, err := variables.GetValueAggregate(variables.TigerdbHost, variables.TigerdbPort, variables.TigerdbUser, variables.TigerdbPass, variables.TigerdbDb)
	if err != nil {
		return nil, fmt.Errorf("unable to create timescale connection: %w", err)
	}

	connUri := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		hostValues[variables.TigerdbUser],
		url.QueryEscape(hostValues[variables.TigerdbPass]),
		hostValues[variables.TigerdbHost],
		hostValues[variables.TigerdbPort],
		hostValues[variables.TigerdbDb],
	)

	config, err := pgxpool.ParseConfig(connUri)
	if err != nil {
		return nil, fmt.Errorf("unable to create timescale connection: %w", err)
	}

	config.MaxConns = 10
	config.MinConns = 5

	conn, err := pgxpool.NewWithConfig(connectionCtx, config)
	if err != nil {
		return nil, fmt.Errorf("unable to create timescale connection: %w", err)
	}

	if err := conn.Ping(connectionCtx); err != nil {
		return nil, fmt.Errorf("unable to create timescale connection: %w", err)
	}

	return &TimescaleConnection{
		conn: conn,
	}, nil
}

func (tc *TimescaleConnection) BeginTx(ctx context.Context) (pgx.Tx, error) {
	return tc.conn.Begin(ctx)
}

func (tc *TimescaleConnection) Pool() *pgxpool.Pool {
	return tc.conn
}

func (tc *TimescaleConnection) Close() {
	tc.conn.Close()
}
