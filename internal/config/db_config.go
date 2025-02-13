package config

type DBConfig struct {
	Host     string
	User     string
	Password string
	DBName   string
	Port     string
	TimeZone string
	SSLMode  string
}

func LoadDBConfigFromEnv() DBConfig {
	return DBConfig{
		Host:     Env("POSTGRES_HOST"),
		User:     Env("POSTGRES_USER"),
		Password: Env("POSTGRES_PASSWORD"),
		DBName:   Env("POSTGRES_DB"),
		Port:     Env("POSTGRES_PORT"),
		TimeZone: Env("APP_TIMEZONE"),
		SSLMode:  Env("POSTGRES_SSL_MODE"),
	}
}
