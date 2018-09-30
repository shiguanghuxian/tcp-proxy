package config

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v1"
)

// Config 配置文件
type Config struct {
	Proxys []*Proxy `yaml:"proxys"`
}

// Proxy 单个代理信息
type Proxy struct {
	Name    string   `yaml:"name"`
	Listen  string   `yaml:"listen"`
	Reverse []string `yaml:"reverse"`
}

// NewConfig 初始化一个server配置文件对象
func NewConfig(path string) (cfg *Config, err error) {
	if path == "" {
		path = "./config/cfg.yaml"
	}
	cfgBytes := make([]byte, 0)
	cfgBytes, err = ioutil.ReadFile(path)
	if err != nil {
		return
	}
	cfg = new(Config)
	err = yaml.Unmarshal(cfgBytes, cfg)
	if err != nil {
		return
	}
	return
}
