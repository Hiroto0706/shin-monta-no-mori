package util

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Environment   string `mapstructure:"ENVIRONMENT"`
	ServerAddress string `mapstructure:"SERVER_ADDRESS"`
	Origin        string `mapstructure:"ORIGIN"`

	// Image
	ImageFetchLimit int `mapstructure:"IMAGE_FETCH_LIMIT"`

	// Character
	CharacterFetchLimit int `mapstructure:"CHARACTER_FETCH_LIMIT"`

	// Category
	CategoryFetchLimit int `mapstructure:"CATEGORY_FETCH_LIMIT"`

	// DB
	DBDriver     string `mapstructure:"DB_DRIVER"`
	DBUser       string `mapstructure:"DB_USER"`
	DBPassword   string `mapstructure:"DB_PASSWORD"`
	DBHost       string `mapstructure:"DB_HOST"`
	DBPort       int    `mapstructure:"DB_PORT"`
	DBName       string `mapstructure:"DB_NAME"`
	DBUrl        string `mapstructure:"DATABASE_URL"`
	TestDBUrl    string `mapstructure:"TEST_DATABASE_URL"`
	TestDBName   string `mapstructure:"TEST_DB_NAME"`
	MigrationURL string `mapstructure:"MIGRATION_URL"`
	SeedURL      string `mapstructure:"SEED_URL"`

	// Redis
	RedisAddress  string `mapstructure:"REDIS_ADDRESS"`
	RedisDbNumber int    `mapstructure:"REDIS_DBNUMBER"`
	RedisPassword string `mapstructure:"REDIS_PASSWORD"`

	// CloudStorage
	BucketName         string `mapstructure:"BUCKET_NAME"`
	CredentialFilePath string `mapstructure:"CREDENTIAL_FILE_PATH"`

	// OTHERS
	TokenSymmetricKey    string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration  time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}

func MakeDBSource(username string, password string, host string, port int, dbname string) string {
	source := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=disable", username, password, host, port, dbname)
	return source
}
