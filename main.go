package main

import (
	"fmt"

	"github.com/conformitea-co/conformitea/cmd"

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

	cmd.Execute()
}
