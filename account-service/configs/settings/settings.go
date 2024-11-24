package settings

import (
	"os"

	"github.com/spf13/viper"
)

func Setup() {
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")

	if _, err := os.Stat(".env"); err == nil {
		if err := viper.ReadInConfig(); err != nil {
			panic(err)
		}
	}

	viper.AutomaticEnv()
}
