package configs

import (
	"github.com/dw-parameter-service/pkg/xlogger"
	"github.com/spf13/viper"
	"log"
	"strings"
)

type ServerConfig struct {
	Port string `mapstructure:"port"`
}

type MongoConfig struct {
	Uri    string `mapstructure:"uri"`
	DBName string `mapstructure:"dbName"`
}

type DBConfig struct {
	Mongo MongoConfig `mapstructure:"mongodb"`
}

type AppConfig struct {
	AppName   string `mapstructure:"appName"`
	DebugMode bool   `mapstructure:"debugMode"`
	// os | file
	LogOutput          string       `mapstructure:"logOutput"`
	LogPath            string       `mapstructure:"logPath"`
	VerboseAPIResponse bool         `mapstructure:"verboseApiResponse"`
	APIServer          ServerConfig `mapstructure:"server"`
	Database           DBConfig     `mapstructure:"database"`
}

var MainConfig AppConfig

func Initialize() error {

	viper.SetConfigType("json")
	viper.AddConfigPath("./")
	viper.SetConfigName("config")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	err := viper.Unmarshal(&MainConfig)
	if err != nil {
		return err
	}

	err = (&xlogger.AppLogger{
		LogPath:     MainConfig.LogPath,
		CompressLog: true,
		DailyRotate: true,
	}).SetAppLogger()

	if err != nil {
		log.Fatalln("Logger error: ", err.Error())
		return err
	}

	xlogger.Log.SetPrefix("[INIT-APP] ")
	xlogger.Log.Println(strings.Repeat("-", 40))
	xlogger.Log.Println("| configuration >> loaded")
	return nil
}
