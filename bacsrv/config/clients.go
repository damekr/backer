package config

import (
    "github.com/spf13/viper"
	"fmt"
    log "github.com/Sirupsen/logrus"
    "os"
)

type Client struct{
    Name    string
    IP      string
}

func init(){
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)

}

func InitClientsConfig(){
     viper.SetConfigName("config")
     viper.AddConfigPath("/home/damian/dev/go/src/github.com/backer/clients")
     err := viper.ReadInConfig()
     if err != nil{
        log.Error("Cannot read clients config file")
     }
     PrintValues()
}


func PrintValues(){
    fmt.Println(viper.AllSettings())
}

