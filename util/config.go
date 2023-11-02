package util

import (
	"github.com/spf13/viper"
)

type ReverseProxyConfig struct {
	Host     string `mapstructure:"host"`
	Endpoint string `mapstructure:"endpoint"`
	Target   string `mapstructure:"target"`
}

type LoadBalancerConfig struct {
	Host     string   `mapstructure:"host"`
	Endpoint string   `mapstructure:"endpoint"`
	Targets  []string `mapstructure:"targets"`
}

type Config struct {
	RPs []ReverseProxyConfig `mapstructure:"reverse-proxies"`
	LBs []LoadBalancerConfig `mapstructure:"load-balancers"`
}

var vp *viper.Viper

func LoadConfig() (Config, error) {
	vp = viper.New()
	var config Config
	vp.SetConfigName("config")
	vp.SetConfigType("json")
	vp.AddConfigPath("./util")
	vp.AddConfigPath(".")
	err := vp.ReadInConfig()
	if err != nil {
		return Config{}, err
	}
	err = vp.Unmarshal(&config)
	if err != nil {
		return Config{}, err
	}
	return config, nil

}
