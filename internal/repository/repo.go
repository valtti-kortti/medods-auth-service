package repository

import (
	"context"
	"fmt"

	"auth-service/internal/config"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

var ErrorNotFound = errors.New("Not found")

type repository struct {
	pool *pgxpool.Pool
}

type Repository interface {
	SetSession(ctx context.Context, session *Session) (int, error)
	DeleteSession(ctx context.Context, id int) error
	GetSession(ctx context.Context, id int) (*Session, error)
}

// инициализация бд
func NewRepository(ctx context.Context, cfg config.Postgres) (Repository, error) {
	connString := fmt.Sprintf(`user=%s password=%s host=%s port=%s dbname=%s sslmode=%s 
        pool_max_conns=%d pool_max_conn_lifetime=%s pool_max_conn_idle_time=%s`,
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
		cfg.SSLMode,
		cfg.PoolMaxConns,
		cfg.PoolMaxConnLifetime,
		cfg.PoolMaxConnIdleTime,
	)

	conf, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse PostgreSQL config")
	}

	conf.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeCacheDescribe

	pool, err := pgxpool.NewWithConfig(ctx, conf)
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to PostgreSQL")
	}

	return &repository{pool: pool}, nil
}

// создание сессии
func (r *repository) SetSession(ctx context.Context, session *Session) (int, error) {
	var id int
	err := r.pool.QueryRow(ctx,
		insertSessionQuery,
		session.UserGUID,
		session.UserAgent,
		session.TokenHash,
		session.IpAddress,
		session.ExpiresAt,
	).Scan(&id)
	if err != nil {
		return 0, errors.Wrap(err, "failed to set refresh token")
	}
	return id, nil
}

// удаление сессии
func (r *repository) DeleteSession(ctx context.Context, id int) error {
	tag, err := r.pool.Exec(ctx, deleteSessionQuery, id)
	if err != nil {
		return errors.Wrap(err, "failed to delete session")
	}

	if tag.RowsAffected() == 0 {
		return ErrorNotFound
	}

	return nil
}

// получение сессии
func (r *repository) GetSession(ctx context.Context, id int) (*Session, error) {
	var session Session
	err := r.pool.QueryRow(ctx, selectSessionQuery, id).Scan(
		&session.SessionID,
		&session.UserGUID,
		&session.UserAgent,
		&session.TokenHash,
		&session.IpAddress,
		&session.ExpiresAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrorNotFound
		}
		return nil, errors.Wrap(err, "failed to get session")
	}
	return &session, nil
}
