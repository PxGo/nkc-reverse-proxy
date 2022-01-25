package modules

import (
	"crypto/tls"
	"net/http/httputil"
)

type NKCHandle struct {
	IsHTTPS      bool
	ReverseProxy *httputil.ReverseProxy
}

type Server struct {
	Id           string   `yaml:"id"`
	Listen       uint16   `yaml:"listen"`
	Name         []string `yaml:"name"`
	SSLKey       string   `yaml:"SSLKey"`
	SSLCert      string   `yaml:"SSLCert"`
	WEBPass      []string `yaml:"WEBPass"`
	WSPass       []string `yaml:"WSPass"`
	WEBType      string   `yaml:"WEBType"`
	WSType       string   `yaml:"WSType"`
	RedirectCode int      `yaml:"redirectCode"`
	RedirectUrl  string   `yaml:"redirectUrl"`
}

type Configs struct {
	Servers []Server `yaml:"servers"`
	Debug   bool     `yaml:"debug"`
}

type ServerPort struct {
	Port      uint16
	TLSConfig *tls.Config
}

type ProxyPass struct {
	WEBPass  []string
	WSPass   []string
	WEBType  string
	WSType   string
	Redirect RedirectInfo
}

type RedirectInfo struct {
	Code int
	Url  string
}
