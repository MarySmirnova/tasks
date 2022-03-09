package config

type Postgres struct {
	Host     string `env:"PG_HOST"`
	Database string `env:"PG_DATABASE"`
	User     string `env:"PG_USER"`
	Password string `env:"PG_PASS"`
	Port     int    `env:"PG_PORT"`
}
