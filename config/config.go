package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config contém todas as configurações da aplicação
type Config struct {
	App        AppConfig
	Database   DatabaseConfig
	JWT        JWTConfig
	Queue      QueueConfig
	Blockchain BlockchainConfig
	Logger     LoggerConfig
}

// AppConfig contém configurações da aplicação
type AppConfig struct {
	Name    string
	Version string
	Env     string
	Port    string
}

// DatabaseConfig contém configurações do banco de dados
type DatabaseConfig struct {
	Driver           string
	ConnectionString string
	MaxOpenConns     int
	MaxIdleConns     int
	MaxConnLifetime  time.Duration
}

// JWTConfig contém configurações de JWT
type JWTConfig struct {
	SecretKey string
	ExpiresIn time.Duration
}

// QueueConfig contém configurações da fila
type QueueConfig struct {
	Type                string // inmemory, rabbitmq, redis
	ConnectionString    string
	QueueName           string
	MaxRetries          int
	RetryBackoffSeconds int
}

// BlockchainConfig contém configurações do blockchain
type BlockchainConfig struct {
	Type            string // simple, ethereum
	ConnectionURL   string
	PrivateKey      string
	ContractAddress string
}

// LoggerConfig contém configurações do logger
type LoggerConfig struct {
	Level string // debug, info, warn, error
}

// Load carrega a configuração do arquivo .env
func Load() (*Config, error) {
	// Carregar arquivo .env se existir
	_ = godotenv.Load()

	return &Config{
		App: AppConfig{
			Name:    getEnv("APP_NAME", "Auditoria Distribuída"),
			Version: getEnv("APP_VERSION", "1.0.0"),
			Env:     getEnv("APP_ENV", "development"),
			Port:    getEnv("APP_PORT", ":8080"),
		},
		Database: DatabaseConfig{
			Driver:           getEnv("DB_DRIVER", "sqlite3"),
			ConnectionString: getEnv("DB_CONNECTION_STRING", "transactions.db"),
			MaxOpenConns:     getEnvInt("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns:     getEnvInt("DB_MAX_IDLE_CONNS", 5),
			MaxConnLifetime:  time.Duration(getEnvInt("DB_MAX_CONN_LIFETIME", 5)) * time.Minute,
		},
		JWT: JWTConfig{
			SecretKey: getEnv("JWT_SECRET_KEY", "your-secret-key-change-in-production"),
			ExpiresIn: time.Duration(getEnvInt("JWT_EXPIRES_IN", 24)) * time.Hour,
		},
		Queue: QueueConfig{
			Type:                getEnv("QUEUE_TYPE", "inmemory"),
			ConnectionString:    getEnv("QUEUE_CONNECTION_STRING", "amqp://guest:guest@localhost:5672/"),
			QueueName:           getEnv("QUEUE_NAME", "transactions"),
			MaxRetries:          getEnvInt("QUEUE_MAX_RETRIES", 3),
			RetryBackoffSeconds: getEnvInt("QUEUE_RETRY_BACKOFF_SECONDS", 60),
		},
		Blockchain: BlockchainConfig{
			Type:            getEnv("BLOCKCHAIN_TYPE", "simple"),
			ConnectionURL:   getEnv("BLOCKCHAIN_CONNECTION_URL", ""),
			PrivateKey:      getEnv("BLOCKCHAIN_PRIVATE_KEY", ""),
			ContractAddress: getEnv("BLOCKCHAIN_CONTRACT_ADDRESS", ""),
		},
		Logger: LoggerConfig{
			Level: getEnv("LOG_LEVEL", "info"),
		},
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

// Validate valida a configuração
func (c *Config) Validate() error {
	if c.App.Port == "" {
		return fmt.Errorf("APP_PORT is required")
	}
	if c.Database.Driver == "" {
		return fmt.Errorf("DB_DRIVER is required")
	}
	if c.JWT.SecretKey == "" || c.JWT.SecretKey == "your-secret-key-change-in-production" {
		return fmt.Errorf("JWT_SECRET_KEY must be set and changed from default value")
	}
	return nil
}
