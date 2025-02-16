package config

type RedisConfig struct {
	Host     string
	Password string
	DBName   string
	Port     string
}

func LoadRedisConfigFromEnv() DBConfig {
	return DBConfig{
		Host:     Env("REDIS_HOST"),
		Port:     Env("REDIS_PORT"),
		Password: Env("REDIS_PASSWORD"),
		DBName:   Env("REDIS_DB"),
	}
}
