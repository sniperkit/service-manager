package env

import (
	"github.com/Peripli/service-manager/logger"
	"github.com/spf13/viper"
)

func Load(location, filename string) error {
	viper.SetConfigName(filename)
	viper.AddConfigPath(location)
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	// this should happen somewhere else
	initializeLogging()
	return nil
}

func initializeLogging() {
	cfg := logger.Configuration{}
	viper.UnmarshalKey("log", &cfg)
	logger.Init(cfg)
}
