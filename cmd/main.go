package main

import (
	"fmt"

	"conformitea/cmd/internal/commands"
	"conformitea/server/config"

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

	var cfg config.Config
	if err := viper.Unmarshal(&cfg); err != nil {
		panic(fmt.Errorf("failed to unmarshal config: %s", err))
	}

	commands.Execute(cfg)
}
