package infra

import "github.com/spf13/viper"

func InitConfig() (err error) {

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	err = viper.ReadInConfig()
	return
}
