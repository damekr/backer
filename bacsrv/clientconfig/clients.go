package clientconfig

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

// DoesClientIntegrated checks if client is in configuration files and is integrated
func DoesClientIntegrated(name string) {
	nil

}

// DoesClientExist checks if client has been added to clients configuration file
func DoesClientExist(name string) {
	nil
}

func InitClientsConfig(c *ServerConfig) {
	log.Info("Client config path: ", c.ClientsConfig)
	//TODO Add checking file
	//TODO Make another instance of Viper, now everything is in one variable, server etc
	clientConfig := viper.New()
	clientConfig.SetConfigFile(c.ClientsConfig)
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
