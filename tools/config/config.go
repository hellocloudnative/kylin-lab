package config

import (
	"fmt"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

var cfgDatabase *viper.Viper
var cfgApplication *viper.Viper
var cfgJwt *viper.Viper
var cfgLog *viper.Viper
var cfgSsl *viper.Viper
var cfkylinCloud *viper.Viper

// 载入配置文件
func ConfigSetup(path string) {
	viper.SetConfigFile(path)
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(fmt.Sprintf("Read config file fail: %s", err.Error()))
	}

	//Replace environment variables
	err = viper.ReadConfig(strings.NewReader(os.ExpandEnv(string(content))))
	if err != nil {
		log.Fatal(fmt.Sprintf("Parse config file fail: %s", err.Error()))
	}

	cfgDatabase = viper.Sub("settings.database")
	if cfgDatabase == nil {
		panic("config not found settings.database")
	}
	DatabaseConfig = InitDatabase(cfgDatabase)

	cfkylinCloud = viper.Sub("settings.kylinCloud")
	if cfkylinCloud == nil {
		panic("config not found settings.kylinCloud")
	}
	KylinCloudConfig = InitKylinCloud(cfkylinCloud)

	cfgApplication = viper.Sub("settings.application")
	if cfgApplication == nil {
		panic("config not found settings.application")
	}
	ApplicationConfig = InitApplication(cfgApplication)

	cfgLog = viper.Sub("settings.log")
	if cfgLog == nil {
		panic("config not found settings.log")
	}
	LogConfig = InitLog(cfgLog)

	cfgSsl = viper.Sub("settings.ssl")
	if cfgSsl == nil {
		panic("config not found settings.ssl")
	}
	SslConfig = InitSsl(cfgSsl)
}

func SetConfig(configPath string, key string, value interface{}) {
	viper.AddConfigPath(configPath)
	viper.Set(key, value)
	viper.WriteConfig()
}
