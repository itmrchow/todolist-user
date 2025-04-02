package infra

import (
	"github.com/spf13/viper"
)

func InitConfig() (err error) {

	viper.AutomaticEnv()
	viper.SetEnvPrefix("APP")

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	err = viper.ReadInConfig()

	return
}
