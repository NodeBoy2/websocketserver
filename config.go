package main

import (
	"log"

	"github.com/spf13/viper"
)

const (
	porxyIP    = "porxyIP"
	porxyPort  = "porxyPort"
	listenIP   = "listenIP"
	listenPort = "listenPort"
)

var configServerDefault = map[string]interface{}{
	porxyIP:    "127.0.0.1",
	porxyPort:  10554,
	listenIP:   "0.0.0.0",
	listenPort: 8080,
}

var configServer *viper.Viper

func initConfig() {
	configServer = viper.New()

	for k, v := range configServerDefault {
		configServer.SetDefault(k, v)
	}
}

func ReadConfig(path, filename string) {
	configServer.SetConfigFile(filename)
	configServer.AddConfigPath(path)
	err := configServer.ReadInConfig()
	if err != nil {
		log.Println(err.Error())
	}

	configServer.WriteConfig()
}

func GetPorxyIP() string {
	return configServer.GetString(porxyIP)
}

func GetPorxyPort() int {
	return configServer.GetInt(porxyPort)
}

func GetListenIP() string {
	return configServer.GetString(listenIP)
}

func GetListenPort() int {
	return configServer.GetInt(listenPort)
}
