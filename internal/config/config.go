package config

import "time"

const EnvPath = ".env"

type AppConfig struct {
	LogLevel string
	Rest     Rest
	Postgres Postgres
	Token    Token
}

type Rest struct {
	ListenAddress string `envconfig:"PORT" required:"true"`
}

type Postgres struct {
	Host                string        `envconfig:"DB_HOST" required:"true"`
	Port                string        `envconfig:"DB_PORT" required:"true"`
	Name                string        `envconfig:"DB_NAME" required:"true"`
	User                string        `envconfig:"DB_USER" required:"true"`
	Password            string        `envconfig:"DB_PASSWORD" required:"true"`
	SSLMode             string        `envconfig:"SSL_MODE" default:"disable"`
	PoolMaxConns        int           `envconfig:"DB_POOL_MAX_CONNS" default:"5"`
	PoolMaxConnLifetime time.Duration `envconfig:"DB_POOL_MAX_CONN_LIFETIME" default:"180s"`
	PoolMaxConnIdleTime time.Duration `envconfig:"DB_POOL_MAX_CONN_IDLE_TIME" default:"100s"`
}

type Token struct {
	Secret   string `envconfig:"JWT_SECRET" required:"true"`
	LenToken int    `envconfig:"LEN_TOKEN" required:"true"`
}
