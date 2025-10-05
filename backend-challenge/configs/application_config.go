package configs

import (
	"github.com/joho/godotenv"
	"oolio.com/kart/constants"
	"os"
	"strconv"
	"time"
)

var (
	AppName    string
	Host       string
	Port       int
	ReleaseEnv string
	LogLevel   string
	DBConfig   DatabaseConfig
)

// DatabaseConfig contains the database configuration
type DatabaseConfig struct {
	Host            string
	Port            int
	User            string
	Password        string
	Name            string
	Schema          string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

// CouponConfiguration contains the coupon configuration
type CouponConfiguration struct {
	Source         string
	DataDir        string
	S3BaseURL      string
	ForceMigration bool
}

// InitApplicationConfig loads the application config from the environment variables or env file
func InitApplicationConfig() error {
	isLocal, _ := strconv.ParseBool(os.Getenv(constants.IsLocal))
	if isLocal {
		err := godotenv.Load(".env")
		if err != nil {
			return err
		}
	}

	return setApplicationConfig()
}

// setApplicationConfig sets the application config from the loaded environment variables
func setApplicationConfig() error {
	var err error

	AppName = os.Getenv(constants.AppName)
	Host = os.Getenv(constants.Host)
	Port, err = strconv.Atoi(os.Getenv(constants.Port))
	if err != nil {
		return err
	}

	// ReleaseEnv can be "local", "dev, "staging", "production"
	ReleaseEnv = os.Getenv(constants.ReleaseEnv)
	if ReleaseEnv == "" {
		ReleaseEnv = "local"
	}

	LogLevel = os.Getenv(constants.LogLevel)

	dbPort, err := strconv.Atoi(getEnvOrDefault(constants.DBPort, "5432"))
	if err != nil {
		return err
	}

	maxOpenConns, err := strconv.Atoi(getEnvOrDefault(constants.DBMaxOpenConns, "20"))
	if err != nil {
		return err
	}

	maxIdleConns, err := strconv.Atoi(getEnvOrDefault(constants.DBMaxIdleConns, "10"))
	if err != nil {
		return err
	}

	connMaxLifetimeSeconds, err := strconv.Atoi(getEnvOrDefault(constants.DBConnMaxLifetime, "1800"))
	if err != nil {
		return err
	}

	DBConfig = DatabaseConfig{
		Host:            os.Getenv(constants.DBHost),
		Port:            dbPort,
		User:            os.Getenv(constants.DBUser),
		Password:        os.Getenv(constants.DBPassword),
		Name:            os.Getenv(constants.DBName),
		Schema:          os.Getenv(constants.DBSchema),
		SSLMode:         os.Getenv(constants.DBSSLMode),
		MaxOpenConns:    maxOpenConns,
		MaxIdleConns:    maxIdleConns,
		ConnMaxLifetime: time.Duration(connMaxLifetimeSeconds) * time.Second,
	}

	return nil
}

// getEnvOrDefault returns the value of the environment variable with the given key, or the fallback value if the environment variable is not set
func getEnvOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
