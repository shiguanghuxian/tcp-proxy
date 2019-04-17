package config

import (
	"io/ioutil"
	"log"

	yaml "gopkg.in/yaml.v1"
)

// Config 配置文件
type Config struct {
	Proxys []*Proxy `yaml:"proxys"`
}

// Proxy 单个代理信息
type Proxy struct {
	Name    string   `yaml:"name"`
	Typ     string   `yaml:"type"`
	Listen  string   `yaml:"listen"`
	Reverse []string `yaml:"reverse"`
}

// NewConfig 初始化一个server配置文件对象
func NewConfig(path string) (cfgChan chan *Config, err error) {
	if path == "" {
		path = "./config/cfg.yaml"
	}
	cfgChan = make(chan *Config, 0)
	// 读取配置文件
	cfg, err := readConfFile(path)
	if err != nil {
		return
	}
	go watcher(cfgChan, path)
	log.Println(1111)
	go func() {
		cfgChan <- cfg
	}()
	log.Println(2222)
	return
}

// ReadConfFile 读取配置文件
func readConfFile(path string) (cfg *Config, err error) {
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

// GetProxyByName 根据name获取一个代理配置
func (c *Config) GetProxyByName(name string) *Proxy {
	for _, v := range c.Proxys {
		if v.Name == name {
			return v
		}
	}
	return nil
}
