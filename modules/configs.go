package modules

import (
	"github.com/go-yaml/yaml"
	"io/ioutil"
	"os"
	"path"
)

var GlobalConfigs *Configs

type Configs struct {
	Servers    []Server `yaml:"servers"`
	ReqLimit   []string `yaml:"req_limit"`
	Console    Console  `yaml:"console"`
	Proxy      bool     `yaml:"proxy"`
	MaxIpCount int16    `yaml:"maxIpCount"`
}

type Console struct {
	Debug   bool `yaml:"debug"`
	Warning bool `yaml:"warning"`
	Error   bool `yaml:"error"`
	Info    bool `yaml:"info"`
}

type Server struct {
	Listen   uint16     `yaml:"listen"`
	Name     []string   `yaml:"name"`
	SSLKey   string     `yaml:"ssl_key"`
	SSLCert  string     `yaml:"ssl_cert"`
	ReqLimit []string   `yaml:"req_limit"`
	Location []Location `yaml:"location"`
}

type Location struct {
	Reg          string   `yaml:"reg"`
	Pass         []string `yaml:"pass"`
	Balance      string   `yaml:"balance"`
	ReqLimit     []string `yaml:"req_limit"`
	RedirectCode int      `yaml:"redirect_code"`
	RedirectUrl  string   `yaml:"redirect_url"`
}

func InitGlobalConfigs() error {
	var err error
	GlobalConfigs, err = GetConfigs()
	if err != nil {
		return err
	}
	return nil
}

func GetConfigsPath() (string, string, error) {
	filePath := "configs.yaml"
	root, err := os.Getwd()
	if err != nil {
		return "", "", err
	}
	if len(os.Args) > 1 {
		filePath = os.Args[1]
	}
	if !path.IsAbs(filePath) {
		filePath = path.Join(root, filePath)
	}
	templateFilePath := path.Join(root, "configs.template.yaml")
	return filePath, templateFilePath, nil
}

func GetConfigs() (*Configs, error) {
	var configs *Configs
	configFilePath, templateConfigFilePath, err := GetConfigsPath()
	if err != nil {
		return nil, err
	}
	file, err := os.ReadFile(configFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			templateFile, err := ioutil.ReadFile(templateConfigFilePath)
			if err != nil {
				return nil, err
			}
			err = ioutil.WriteFile(configFilePath, templateFile, 0644)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
		file, err = os.ReadFile(configFilePath)
		if err != nil {
			return nil, err
		}
	}
	err = yaml.Unmarshal(file, &configs)
	if err != nil {
		return nil, err
	}
	return configs, nil
}
