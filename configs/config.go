package configs

import (
	"fmt"

	"github.com/spf13/viper"
)

// LoadConfig loads env variables
func LoadConfig() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s", err.Error()))
	}
}
