package tools

import (
	"errors"
	"github.com/go-yaml/yaml"
	"os"
)

type Transport struct {
	KeepAlive                 bool  `yaml:"keep-alive"`
	MaxIdleConnections        int64 `yaml:"maxIdleConnections"`
	MaxIdleConnectionsPerHost int64 `yaml:"MaxIdleConnectionsPerHost"`
	MaxConnectionsPerHost     int64 `yaml:"maxConnectionsPerHost"`
}

type Server struct {
	Id           string   `yaml:"id"`
	Listen       int64    `yaml:"listen"`
	Name         []string `yaml:"name"`
	SSLKey       string   `yaml:"SSLKey"`
	SSLCert      string   `yaml:"SSLCert"`
	WEBPass      []string `yaml:"WEBPass"`
	WSPass       []string `yaml:"WSPass"`
	WEBType      string   `yaml:"WEBType"`
	WSType       string   `yaml:"WSType"`
	RedirectCode int64    `yaml:"redirectCode"`
	RedirectUrl  string   `yaml:"redirectUrl"`
}

type Configs struct {
	Transport Transport `yaml:"transport"`
	Servers   []Server  `yaml:"servers"`
}

func GetConfigsPath() (string, error) {
	filePath := os.Args[1]
	if len(filePath) == 0 {
		return "", errors.New("未指定配置文件路径")
	}
	return filePath, nil
}

func GetConfigs() (*Configs, error) {
	configFilePath, err := GetConfigsPath()
	if err != nil {
		return nil, err
	}
	file, err := os.ReadFile(configFilePath)
	if err != nil {
		return nil, err
	}
	var configs Configs
	err = yaml.Unmarshal(file, &configs)
	if err != nil {
		return nil, err
	}
	return &configs, nil
}

func GetPortFromConfigs() ([]int64, error) {
	configs, err := GetConfigs()
	if err != nil {
		return nil, err
	}
	var port []int64
	for _, server := range configs.Servers {
		port = append(port, server.Listen)
	}
	return port, nil
}
