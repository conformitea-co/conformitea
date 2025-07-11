package sections

import (
	"fmt"
)

type RedisConfig struct {
	Address  string `mapstructure:"address"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
}

func (r *RedisConfig) Validate() error {
	if r.Address == "" {
		return fmt.Errorf("redis.address is required")
	}

	return nil
}
