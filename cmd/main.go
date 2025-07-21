package main

import (
	"fmt"

	"conformitea/cmd/config"
	"conformitea/cmd/internal/commands"

	"github.com/spf13/viper"
)

func main() {
	viper.SetEnvPrefix("CONFORMITEA")
	viper.AutomaticEnv()

	viper.SetConfigType("toml")
	viper.SetConfigName("conformitea")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			panic(fmt.Errorf("fatal error config file: %s", err))
		}
	}

	var config config.Config
	if err := viper.Unmarshal(&config); err != nil {
		panic(fmt.Errorf("failed to unmarshal config: %s", err))
	}

	commands.Execute(config)
}
