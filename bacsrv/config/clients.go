package config

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
)

type Client struct {
	Name string
	IP   string
}

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)

}

func InitClientsConfig(c *ServerConfig) {
	log.Info("Client config path: ", c.ClientsConfig)
	//TODO Add checking file
	//TODO Make another instance of Viper, now everything is in one variable, server etc
	viper.SetConfigFile(c.ClientsConfig)
	//viper.SetConfigName("config")
	//viper.AddConfigPath("/home/damian/dev/go/src/github.com/backer/clients")
	//err := viper.ReadInConfig()
	//if err != nil{
	//   log.Error("Cannot read clients config file")
	//}
	//PrintValues()
}

func PrintValues() {
	fmt.Println(viper.AllSettings())
}
