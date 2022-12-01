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
	Id              string   `yaml:"id"`
	Listen          uint16   `yaml:"listen"`
	Name            []string `yaml:"name"`
	SSLKey          string   `yaml:"ssl_key"`
	SSLCert         string   `yaml:"ssl_cert"`
	Pass            []string `yaml:"pass"`
	SocketIoPass    []string `yaml:"socket_io_pass"`
	Balance         string   `yaml:"balance"`
	SocketIoBalance string   `yaml:"socket_io_balance"`
	RedirectCode    int      `yaml:"redirect_code"`
	RedirectUrl     string   `yaml:"redirect_url"`
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
