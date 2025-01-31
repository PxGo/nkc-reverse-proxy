package modules

import (
	_ "embed"
	"errors"
	"flag"
	"fmt"
	"gopkg.in/yaml.v2"
	"nkc-reverse-proxy/conf"
	"os"
	"path/filepath"
)

var GlobalConfigs *Configs

var InnerConfigs *Configs

type Configs struct {
	Servers    []Server `yaml:"servers"`
	ReqLimit   []string `yaml:"req_limit"`
	IpLimit    []string `yaml:"ip_limit"`
	Console    Console  `yaml:"console"`
	Proxy      bool     `yaml:"proxy"`
	MaxIpCount int16    `yaml:"maxIpCount"`
	Template   Template `yaml:"template"`
}

type TemplateContent struct {
	Title string `yaml:"title"`
	Desc  string `yaml:"desc"`
}

type Template struct {
	Page404 TemplateContent `yaml:"page404"`
	Page500 TemplateContent `yaml:"page500"`
	Page503 TemplateContent `yaml:"page503"`
	Page429 TemplateContent `yaml:"page429"`
	Page403 TemplateContent `yaml:"page403"`
}

type Console struct {
	Debug   bool `yaml:"debug"`
	Warning bool `yaml:"warning"`
	Error   bool `yaml:"error"`
	Info    bool `yaml:"info"`
}

type Server struct {
	Listen   uint16          `yaml:"listen"`
	Name     []string        `yaml:"name"`
	SSLKey   string          `yaml:"ssl_key"`
	SSLCert  string          `yaml:"ssl_cert"`
	SSLAuto  bool            `yaml:"ssl_auto"`
	ReqLimit []string        `yaml:"req_limit"`
	IpLimit  []string        `yaml:"ip_limit"`
	Location []Location      `yaml:"location"`
	Page404  TemplateContent `yaml:"page404"`
	Page500  TemplateContent `yaml:"page500"`
	Page503  TemplateContent `yaml:"page503"`
	Page429  TemplateContent `yaml:"page429"`
	Page403  TemplateContent `yaml:"page403"`
}

type Location struct {
	Reg          string   `yaml:"reg"`
	Pass         []string `yaml:"pass"`
	Balance      string   `yaml:"balance"`
	ReqLimit     []string `yaml:"req_limit"`
	IpLimit      []string `yaml:"ip_limit"`
	RedirectCode int      `yaml:"redirect_code"`
	RedirectUrl  string   `yaml:"redirect_url"`
	Root         string   `yaml:"root"`
	RootPrefix   string   `yaml:"rootPrefix"`
}

func InitGlobalConfigs() error {
	var err error
	GlobalConfigs, err = GetConfigs()
	if err != nil {
		return err
	}
	return nil
}

func GetConfigsPath() (string, error) {
	filename := "config.yaml"
	root, err := os.Getwd()
	if err != nil {
		return "", err
	}

	defaultFilePath := filepath.Join(root, filename)

	var filePath string
	var help bool
	var version bool

	flag.StringVar(&filePath, "f", defaultFilePath, "Path to the configuration file")
	flag.BoolVar(&help, "h", false, "show help information")
	flag.BoolVar(&version, "v", false, "show version")

	flag.Parse()

	if help {
		fmt.Println("Usage: NRP [options]")
		flag.PrintDefaults()
		os.Exit(0)
	}

	if version {
		fmt.Println(fmt.Sprintf("NRP %s", conf.CodeVersion))
		os.Exit(0)
	}

	if !filepath.IsAbs(filePath) {
		filePath = filepath.Join(root, filePath)
	}

	logger.InfoLog(fmt.Sprintf("Configuration file path is %s", filePath))

	return filePath, nil
}

func GetConfigs() (*Configs, error) {

	if InnerConfigs != nil {
		return InnerConfigs, nil
	}
	var configs *Configs
	configFilePath, err := GetConfigsPath()
	if err != nil {
		return nil, err
	}
	file, err := os.ReadFile(configFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, errors.New(fmt.Sprintf("Configuration file not found at %s", configFilePath))
		} else {
			return nil, err
		}
	}
	err = yaml.Unmarshal(file, &configs)
	if err != nil {
		return nil, err
	}

	InnerConfigs = configs

	return InnerConfigs, nil
}
