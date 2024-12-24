package infra

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var Config *viper.Viper

// List of allowed environments this application expects to receive.
var allowedEnvironments = []string{
	"test",
	"dev",
	"hml",
	"prd",
}

var (
	ErrInvalidEnvironment = func(env string) error {
		return fmt.Errorf("'%s' is not a valid environment. Use one of %+v", env, allowedEnvironments)
	}
)

func init() {
	Config = viper.New()

	Config.AddConfigPath("./")
	Config.SetConfigFile(".env")

	err := Config.ReadInConfig()
	if err != nil {
		logrus.WithError(err).Fatal("Failed to read configuration")
	}

}
