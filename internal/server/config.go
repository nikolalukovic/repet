package server

import (
	"github.com/spf13/viper"
)

type RepetConfig struct {
	Port    int32
	Address string
}

func ConsumeConfiguration() (*RepetConfig, error) {
	viper.SetEnvPrefix("REPET")
	viper.AutomaticEnv()

	if !viper.IsSet("ADDR") {
		return nil, &RepetError{Code: ConfigAddrNotSet}
	}

	if !viper.IsSet("PORT") {
		return nil, &RepetError{Code: ConfigPortNotSet}
	}

	addr := viper.GetString("ADDR")
	port := viper.GetInt32("PORT")

	return &RepetConfig{
		Address: addr,
		Port:    port,
	}, nil
}
