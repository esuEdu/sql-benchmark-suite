package db

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresDB struct {
	URI   string
	Pool  *pgxpool.Pool
	Table string
}

func NewPostgres(uri, table string) (*PostgresDB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, uri)
	if err != nil {
		return nil, err
	}

	createTable := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			index INT PRIMARY KEY,
			value TEXT,
			ts BIGINT
		);
	`, table)

	if _, err := pool.Exec(ctx, createTable); err != nil {
		return nil, err
	}

	return &PostgresDB{
		URI:   uri,
		Pool:  pool,
		Table: table,
	}, nil
}

func (p *PostgresDB) Name() string { return "postgres" }

func (p *PostgresDB) WriteTest(n int) (time.Duration, error) {
	ctx := context.Background()
	start := time.Now()

	query := fmt.Sprintf(`INSERT INTO %s (index, value, ts) VALUES ($1, $2, $3)
		ON CONFLICT (index) DO UPDATE SET value = EXCLUDED.value, ts = EXCLUDED.ts`, p.Table)

	for i := 0; i < n; i++ {
		_, err := p.Pool.Exec(ctx, query,
			i,
			fmt.Sprintf("value-%d", i),
			time.Now().UnixNano(),
		)
		if err != nil {
			return 0, err
		}
	}

	return time.Since(start), nil
}

func (p *PostgresDB) ReadTest(n int) (time.Duration, error) {
	ctx := context.Background()
	start := time.Now()

	query := fmt.Sprintf(`SELECT * FROM %s `, p.Table)

	_ = p.Pool.QueryRow(ctx, query).Scan(new(string))

	return time.Since(start), nil
}
