package config

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

const (
	httpPort string = "HTTP_PORT"
	// Log configuration
	loggerJSON  string = "LOGGING_JSON"
	loggerLevel string = "LOGGING_LEVEL"

	// Postgres config
	PgHost       string = "POSTGRES_HOST_ADDRESS"
	PgUser       string = "POSTGRES_USER"
	PgPassword   string = "POSTGRES_PASSWORD"
	PgDatabase   string = "POSTGRES_DBNAME"
	PgPort       string = "POSTGRES_PORT"
	PgSSLMode    string = "POSTGRES_SSL_MODE"
	PgLogQueries string = "POSTGRES_LOG"
)

// AppConfig The application configuration definition
type AppConfig struct {
	HTTPPort       int
	LoggerConfig   LoggerConfig
	PostgresConfig PostgresConfig
}

type LoggerConfig struct {
	Level zerolog.Level
	JSON  bool
}
type PostgresConfig struct {
	Host       string
	User       string
	Password   string `json:"-"`
	Database   string
	Port       int
	SSLMode    string
	LogQueries bool
}

func NewAppConfig() AppConfig {
	// Load config file
	viper.SetConfigFile(".env")

	// Read the config file
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("Error reading config file: %s\n", err)
	}

	// Automatic environment variables
	viper.AutomaticEnv()

	return AppConfig{
		HTTPPort:       viper.GetInt(httpPort),
		LoggerConfig:   NewLoggerConfig(),
		PostgresConfig: NewPostgresConfig(),
	}
}

func NewPostgresConfig() PostgresConfig {
	pgConfig := PostgresConfig{
		Host:       viper.GetString(PgHost),
		User:       viper.GetString(PgUser),
		Password:   viper.GetString(PgPassword),
		Database:   viper.GetString(PgDatabase),
		Port:       viper.GetInt(PgPort),
		SSLMode:    viper.GetString(PgSSLMode),
		LogQueries: viper.GetBool(PgLogQueries),
	}
	if pgConfig.Host == "" || pgConfig.User == "" || pgConfig.Password == "" || pgConfig.Database == "" ||
		pgConfig.Port == 0 || pgConfig.SSLMode == "" {
		log.Fatal().Msg("Invalid database config parameters")
	}
	return pgConfig
}

func NewLoggerConfig() LoggerConfig {
	loggingLevel, err := zerolog.ParseLevel(viper.GetString(loggerLevel))
	if err != nil || loggingLevel == zerolog.NoLevel {
		loggingLevel = zerolog.InfoLevel
	}
	loggerConfig := LoggerConfig{
		Level: loggingLevel,
		JSON:  viper.GetBool(loggerJSON),
	}
	return loggerConfig
}
