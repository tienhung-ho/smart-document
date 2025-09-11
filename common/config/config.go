package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Environment string         `mapstructure:"environment"`
	Server      ServerConfig   `mapstructure:"server"`
	Database    DatabaseConfig `mapstructure:"database"`
	Redis       RedisConfig    `mapstructure:"redis"`
	Kafka       KafkaConfig    `mapstructure:"kafka"`
	JWT         JWTConfig      `mapstructure:"jwt"`
	Logging     LoggingConfig  `mapstructure:"logging"`
	Storage     StorageConfig  `mapstructure:"storage"`
}

type ServerConfig struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	ReadTimeout  int    `mapstructure:"read_timeout"`
	WriteTimeout int    `mapstructure:"write_timeout"`
}

type DatabaseConfig struct {
	Driver   string `mapstructure:"driver"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
	SSLMode  string `mapstructure:"ssl_mode"`
	MaxIdle  int    `mapstructure:"max_idle"`
	MaxOpen  int    `mapstructure:"max_open"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	Database int    `mapstructure:"database"`
}

type KafkaConfig struct {
	Brokers []string `mapstructure:"brokers"`
	GroupID string   `mapstructure:"group_id"`
}

type JWTConfig struct {
	Secret    string `mapstructure:"secret"`
	ExpiresIn int    `mapstructure:"expires_in"`
}

type LoggingConfig struct {
	Level      string `mapstructure:"level"`
	Format     string `mapstructure:"format"`
	Output     string `mapstructure:"output"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxBackups int    `mapstructure:"max_backups"`
	MaxAge     int    `mapstructure:"max_age"`
	Compress   bool   `mapstructure:"compress"`
}

type StorageConfig struct {
	Type      string      `mapstructure:"type"`
	Minio     MinioConfig `mapstructure:"minio"`
	LocalPath string      `mapstructure:"local_path"`
}

type MinioConfig struct {
	Endpoint        string `mapstructure:"endpoint"`
	AccessKeyID     string `mapstructure:"access_key_id"`
	SecretAccessKey string `mapstructure:"secret_access_key"`
	UseSSL          bool   `mapstructure:"use_ssl"`
	BucketName      string `mapstructure:"bucket_name"`
}

func LoadConfig(configPath string, serviceName string) (*Config, error) {
	v := viper.New()

	// Set default values
	setDefaults(v)

	// Configure Viper
	v.SetConfigName(serviceName)
	v.SetConfigType("yaml")
	v.AddConfigPath(configPath)
	v.AddConfigPath("./etc")
	v.AddConfigPath("../etc")
	v.AddConfigPath("../../etc")

	// Enable environment variable reading
	v.AutomaticEnv()
	v.SetEnvPrefix("SD") // Smart Document prefix
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Read config file
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
		// Config file not found, continue with env vars and defaults
	}

	// Unmarshal config
	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Override with environment-specific values
	if err := loadEnvironmentSpecific(v, &config); err != nil {
		return nil, fmt.Errorf("failed to load environment-specific config: %w", err)
	}

	return &config, nil
}

func setDefaults(v *viper.Viper) {
	// Environment
	v.SetDefault("environment", "development")

	// Server defaults
	v.SetDefault("server.host", "0.0.0.0")
	v.SetDefault("server.port", 8080)
	v.SetDefault("server.read_timeout", 30)
	v.SetDefault("server.write_timeout", 30)

	// Database defaults
	v.SetDefault("database.driver", "postgres")
	v.SetDefault("database.host", "localhost")
	v.SetDefault("database.port", 5432)
	v.SetDefault("database.username", "postgres")
	v.SetDefault("database.password", "postgres")
	v.SetDefault("database.database", "smart_document")
	v.SetDefault("database.ssl_mode", "disable")
	v.SetDefault("database.max_idle", 10)
	v.SetDefault("database.max_open", 100)

	// Redis defaults
	v.SetDefault("redis.host", "localhost")
	v.SetDefault("redis.port", 6379)
	v.SetDefault("redis.password", "")
	v.SetDefault("redis.database", 0)

	// Kafka defaults
	v.SetDefault("kafka.brokers", []string{"localhost:9092"})
	v.SetDefault("kafka.group_id", "smart-document")

	// JWT defaults
	v.SetDefault("jwt.secret", "your-secret-key")
	v.SetDefault("jwt.expires_in", 3600)

	// Logging defaults
	v.SetDefault("logging.level", "info")
	v.SetDefault("logging.format", "json")
	v.SetDefault("logging.output", "stdout")
	v.SetDefault("logging.max_size", 100)
	v.SetDefault("logging.max_backups", 3)
	v.SetDefault("logging.max_age", 28)
	v.SetDefault("logging.compress", true)

	// Storage defaults
	v.SetDefault("storage.type", "minio")
	v.SetDefault("storage.minio.endpoint", "localhost:9000")
	v.SetDefault("storage.minio.access_key_id", "minioadmin")
	v.SetDefault("storage.minio.secret_access_key", "minioadmin")
	v.SetDefault("storage.minio.use_ssl", false)
	v.SetDefault("storage.minio.bucket_name", "documents")
	v.SetDefault("storage.local_path", "./uploads")
}

func loadEnvironmentSpecific(v *viper.Viper, config *Config) error {
	environment := v.GetString("environment")

	// Load environment-specific overrides
	envFile := fmt.Sprintf("config-%s.yaml", environment)

	// Try to read environment-specific config
	envViper := viper.New()
	envViper.SetConfigName(strings.TrimSuffix(envFile, ".yaml"))
	envViper.SetConfigType("yaml")
	envViper.AddConfigPath("./etc")
	envViper.AddConfigPath("../etc")
	envViper.AddConfigPath("../../etc")

	if err := envViper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return fmt.Errorf("failed to read environment config file: %w", err)
		}
		// Environment-specific config not found, which is fine
	} else {
		// Merge environment-specific config
		if err := envViper.Unmarshal(config); err != nil {
			return fmt.Errorf("failed to unmarshal environment config: %w", err)
		}
	}

	return nil
}

func GetEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func MustGetEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic(fmt.Sprintf("required environment variable %s is not set", key))
	}
	return value
}
