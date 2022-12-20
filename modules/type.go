package modules

import (
	"crypto/tls"
	"net/http/httputil"
)

type NKCHandle struct {
	IsHTTPS      bool
	Port         uint16
	ReverseProxy *httputil.ReverseProxy
}

type Configs struct {
	Servers    []Server `yaml:"servers"`
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
	Id       string     `yaml:"id"`
	Listen   uint16     `yaml:"listen"`
	Name     []string   `yaml:"name"`
	SSLKey   string     `yaml:"ssl_key"`
	SSLCert  string     `yaml:"ssl_cert"`
	Location []Location `yaml:"location"`
}

type Location struct {
	Reg          string   `yaml:"reg"`
	Pass         []string `yaml:"pass"`
	Balance      string   `yaml:"balance"`
	ReqLimitIp   []uint16 `yaml:"req_limit_ip"`
	RedirectCode int      `yaml:"redirect_code"`
	RedirectUrl  string   `yaml:"redirect_url"`
}

type ServerPort struct {
	Port      uint16
	TLSConfig *tls.Config
}

type NameLocation map[string][]Location
type ServerLocation map[uint16]NameLocation

type ProxyPass struct {
	Pass            []string
	SocketIoPass    []string
	Balance         string
	SocketIoBalance string
	Redirect        RedirectInfo
}

type RedirectInfo struct {
	Code int
	Url  string
}
