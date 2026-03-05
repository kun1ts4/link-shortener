package postgres

import (
	"context"
	"fmt"
	"link-shortener/internal/config"
	"link-shortener/internal/domain"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PgRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(cfg config.PostgresConfig) (*PgRepository, error) {
	poolConfig, err := pgxpool.ParseConfig(cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	poolConfig.MaxConns = cfg.MaxConnections
	poolConfig.MinConns = cfg.MinConnections

	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = pool.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
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

func (r *PgRepository) Create(link *domain.Link) error {
	dbLink := FromDomain(link)

	query := `INSERT INTO links (original, short, clicks, created_at) VALUES ($1, $2, $3, $4)`
	_, err := r.pool.Exec(context.Background(), query, dbLink.Original, dbLink.Short, dbLink.Clicks, dbLink.CreatedAt)

	if err != nil {
		return fmt.Errorf("create link in db: %w", err)
	}

	return nil
}

func (r *PgRepository) FindByShort(short string) (*domain.Link, error) {
	query := `SELECT id, original, short, clicks, created_at FROM links WHERE short = $1`

	var dbLink LinkDB
	err := r.pool.QueryRow(context.Background(), query, short).Scan(&dbLink.ID, &dbLink.Original, &dbLink.Short, &dbLink.Clicks, &dbLink.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("find by short %s: %w", short, domain.ErrNotFound)
	}

	return dbLink.ToDomain(), nil
}

func (r *PgRepository) FindByOriginal(original string) (*domain.Link, error) {
	query := `SELECT id, original, short, clicks, created_at FROM links WHERE original = $1`

	var dbLink LinkDB
	err := r.pool.QueryRow(context.Background(), query, original).Scan(&dbLink.ID, &dbLink.Original, &dbLink.Short, &dbLink.Clicks, &dbLink.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("find by original: %w", domain.ErrNotFound)
	}

	return dbLink.ToDomain(), nil
}

func (r *PgRepository) FindByID(id int64) (*domain.Link, error) {
	query := `SELECT id, original, short, clicks, created_at FROM links WHERE id = $1`

	var dbLink LinkDB
	err := r.pool.QueryRow(context.Background(), query, id).Scan(&dbLink.ID, &dbLink.Original, &dbLink.Short, &dbLink.Clicks, &dbLink.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("find by id %d: %w", id, domain.ErrNotFound)
	}

	return dbLink.ToDomain(), nil
}

func (r *PgRepository) IncrementClicks(id int64) error {
	query := `UPDATE links SET clicks = clicks + 1 WHERE id = $1`

	_, err := r.pool.Exec(context.Background(), query, id)
	if err != nil {
		return fmt.Errorf("increment clicks for id %d: %w", id, err)
	}

	return nil
}
