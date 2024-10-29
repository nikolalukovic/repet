package server

import (
	"github.com/nikolalukovic/repet/internal/shared"
	"github.com/spf13/viper"
)

type RepetConfig struct {
	Port    int32  `env:"REPET_PORT"`
	Address string `env:"REPET_ADDR"`
}

func ConsumeConfiguration() (*RepetConfig, error) {
	viper.SetEnvPrefix("REPET")
	viper.AutomaticEnv()

	if !viper.IsSet("ADDR") {
		return nil, &shared.RepetError{Code: shared.ConfigAddrNotSet}
	}

	if !viper.IsSet("PORT") {
		return nil, &shared.RepetError{Code: shared.ConfigPortNotSet}
	}

	addr := viper.GetString("ADDR")
	port := viper.GetInt32("PORT")

	return &RepetConfig{
		Address: addr,
		Port:    port,
	}, nil
}
