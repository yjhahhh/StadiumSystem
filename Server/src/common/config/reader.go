package config

import (
	"fmt"
	"sync"

	"github.com/spf13/viper"
)

const (
	ConfigFile = "../webserver/config/globalconfig.toml"
)

// 全局配置
var gConfig *Config
var once sync.Once

func InitGlobalConfig() {
	setDefault()

	var conf Config
	viper.SetConfigFile(ConfigFile)
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	if err := viper.Unmarshal(&conf); err != nil {
        panic(fmt.Errorf("unmarshal config error: %s", err.Error()))
    }

	once.Do(func() {
		gConfig = &conf
	})
	fmt.Printf("globalconfig: %+v\n", conf)
}

func GetGlobalConfig() *Config {
	return gConfig
}

func GetMySQLConfig() *MySQLConfig {
	return &gConfig.MySQLConf
}

func GetRedisConfig() *RedisConfig {
	return &gConfig.RedisConf
}

func setDefault() {
	viper.SetDefault("Addr", "localhost:8080")

}

func GetString(key string) (string, bool) {
	if exists := viper.IsSet(key); !exists {
		return "", false
	}
	return viper.GetString(key), true
}