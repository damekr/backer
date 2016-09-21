package main

import (
    "github.com/spf13/viper"
	"fmt"
)

type ServerConfig struct {
    MgmtPort    string
    DataPort    string
    LogOutput   string // STDOUT, FILE, SYSLOG
    Debug       bool

}

func fillConfigStruct() *ServerConfig{
    return &ServerConfig {
        MgmtPort: viper.GetString("server.MgmtPort"),
        DataPort: viper.GetString("server.DataPort"),
        LogOutput: viper.GetString("server.LogOutput"),
        Debug: viper.GetBool("server.Debug"),
    }
}

func(c *ServerConfig) showConfig(){
    fmt.Printf("Config Struct: %#v\n", c)
}


func setConfigPath(){
    // Viper can cooperate with Cobra arg parser consider reading config file path from arg
    viper.SetConfigName("bacsrv")
    viper.AddConfigPath("/home/damian/dev/go/src/github.com/backer/config")
}


func readConfigFile() *ServerConfig{
    setConfigPath()
    viper.ReadInConfig()
    config := fillConfigStruct()
    return config
}