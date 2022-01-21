package modules

import (
	"crypto/tls"
	"net/http/httputil"
	"time"
)

type NKCHandle struct {
	IsHTTPS      bool
	ReverseProxy *httputil.ReverseProxy
}

type Transport struct {
	KeepAlive                 bool          `yaml:"keeplive"`
	MaxIdleConnections        int           `yaml:"maxIdleConnections"`
	MaxIdleConnectionsPerHost int           `yaml:"MaxIdleConnectionsPerHost"`
	MaxConnectionsPerHost     int           `yaml:"maxConnectionsPerHost"`
	Timeout                   time.Duration `yaml:"timeout"`
	KeepAliveTimeout          time.Duration `yaml:"keepAliveTimeout"`
	IdleConnectionTimeout     time.Duration `yaml:"idleConnectionTimeout"`
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
	RedirectCode int      `yaml:"redirectCode"`
	RedirectUrl  string   `yaml:"redirectUrl"`
}

type Configs struct {
	Transport Transport `yaml:"transport"`
	Servers   []Server  `yaml:"servers"`
	ErrorLog  string    `yaml:"errorLog"`
	PProf     int64     `yaml:"pprof"`
}

type ServerPort struct {
	Port      int64
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
