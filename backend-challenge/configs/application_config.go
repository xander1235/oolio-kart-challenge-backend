package configs

import (
	"github.com/joho/godotenv"
	"oolio.com/kart/constants"
	"os"
	"strconv"
)

var (
	AppName    string
	Host       string
	Port       int
	ReleaseEnv string
	LogLevel   string
)

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
	os.Getenv(constants.IsLocal)

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

	return nil
}
