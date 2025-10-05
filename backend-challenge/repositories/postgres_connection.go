package repositories

import (
	"context"
	"fmt"
	"net/url"

	"github.com/jackc/pgx/v5/pgxpool"

	"oolio.com/kart/configs"
)

var postgresPool *pgxpool.Pool

// Initialize establishes a shared pgx connection pool based on the loaded configuration.
func Initialize(ctx context.Context) error {
	if postgresPool != nil {
		return nil
	}

	connString, err := buildConnectionString()
	if err != nil {
		return err
	}

	cfg, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return err
	}

	if maxConns := configs.DBConfig.MaxOpenConns; maxConns > 0 {
		cfg.MaxConns = int32(maxConns)
	}

	if idleConns := configs.DBConfig.MaxIdleConns; idleConns > 0 {
		cfg.MinConns = int32(idleConns)
	}

	if lifetime := configs.DBConfig.ConnMaxLifetime; lifetime > 0 {
		cfg.MaxConnLifetime = lifetime
	}

	postgresPool, err = pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return err
	}

	return nil
}

// Pool returns the initialized pgx pool instance.
func Pool() (*pgxpool.Pool, error) {
	if postgresPool == nil {
		return nil, fmt.Errorf("postgres pool not initialized")
	}

	return postgresPool, nil
}

// Close terminates the shared connection pool.
func Close() {
	if postgresPool != nil {
		postgresPool.Close()
		postgresPool = nil
	}
}

// buildConnectionString builds a postgres connection string based on the loaded configuration.
func buildConnectionString() (string, error) {
	cfg := configs.DBConfig

	if cfg.Host == "" {
		return "", fmt.Errorf("database host is required")
	}

	if cfg.Name == "" {
		return "", fmt.Errorf("database name is required")
	}

	host := cfg.Host
	if cfg.Port > 0 {
		host = fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	}

	u := &url.URL{
		Scheme: "postgres",
		Host:   host,
		Path:   "/" + cfg.Name,
	}

	if cfg.User != "" {
		if cfg.Password != "" {
			u.User = url.UserPassword(cfg.User, cfg.Password)
		} else {
			u.User = url.User(cfg.User)
		}
	}

	query := u.Query()
	if cfg.SSLMode != "" {
		query.Set("sslmode", cfg.SSLMode)
	}
	u.RawQuery = query.Encode()

	return u.String(), nil
}
