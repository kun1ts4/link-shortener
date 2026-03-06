package postgres

import (
	"context"
	"errors"
	"fmt"
	"link-shortener/internal/config"
	"link-shortener/internal/domain"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PgRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(cfg config.PostgresConfig) (*PgRepository, error) {
	poolConfig, err := pgxpool.ParseConfig(cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	poolConfig.MaxConns = cfg.MaxConnections
	poolConfig.MinConns = cfg.MinConnections
	poolConfig.MaxConnLifetime = 1 * time.Hour
	poolConfig.MaxConnIdleTime = 30 * time.Minute
	poolConfig.HealthCheckPeriod = 1 * time.Minute

	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, fmt.Errorf("create pool: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = pool.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("ping database: %w", err)
	}

	return &PgRepository{
		pool: pool,
	}, nil
}

func (r *PgRepository) Close() {
	if r.pool != nil {
		r.pool.Close()
	}
}

func (r *PgRepository) Create(ctx context.Context, link *domain.Link) error {
	dbLink := FromDomain(link)

	query := `INSERT INTO links (original, short, clicks, created_at) VALUES ($1, $2, $3, $4)`
	_, err := r.pool.Exec(ctx, query, dbLink.Original, dbLink.Short, dbLink.Clicks, dbLink.CreatedAt)
	if err != nil {
		return fmt.Errorf("create link: %w", err)
	}

	return nil
}

func (r *PgRepository) FindByShort(ctx context.Context, short string) (*domain.Link, error) {
	query := `SELECT id, original, short, clicks, created_at FROM links WHERE short = $1`

	var dbLink LinkDB
	err := r.pool.QueryRow(ctx, query, short).Scan(&dbLink.ID, &dbLink.Original, &dbLink.Short, &dbLink.Clicks, &dbLink.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("find by short: %w", err)
	}

	return dbLink.ToDomain(), nil
}

func (r *PgRepository) FindByOriginal(ctx context.Context, original string) (*domain.Link, error) {
	query := `SELECT id, original, short, clicks, created_at FROM links WHERE original = $1`

	var dbLink LinkDB
	err := r.pool.QueryRow(ctx, query, original).Scan(&dbLink.ID, &dbLink.Original, &dbLink.Short, &dbLink.Clicks, &dbLink.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("find by original: %w", err)
	}

	return dbLink.ToDomain(), nil
}

func (r *PgRepository) IncrementClicks(ctx context.Context, short string) error {
	query := `UPDATE links SET clicks = clicks + 1 WHERE short = $1`

	tag, err := r.pool.Exec(ctx, query, short)
	if err != nil {
		return fmt.Errorf("increment clicks for %s: %w", short, err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}

	return nil
}
