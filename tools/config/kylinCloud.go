package config

import "github.com/spf13/viper"

type KylinCloud struct {
	AuthUrl string
	ApiUrl  string
}

func InitKylinCloud(cfg *viper.Viper) *KylinCloud {
	return &KylinCloud{
		AuthUrl: cfg.GetString("authurl"),
		ApiUrl:  cfg.GetString("apiurl"),
	}
}

var KylinCloudConfig = new(KylinCloud)
