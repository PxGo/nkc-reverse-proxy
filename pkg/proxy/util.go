package proxy

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"

	"github.com/tokisakiyuu/nkc-proxy-go-pure/pkg/config"
)

func isWebSocketUpgradeRequest(req *http.Request) bool {
	connHeader := req.Header.Get("Connection")
	if strings.ToLower(connHeader) != "upgrade" {
		return false
	}
	upgradeHeader := req.Header.Get("Upgrade")
	return strings.ToLower(upgradeHeader) == "websocket"
}

func isSocketIOPolling(req *http.Request) bool {
	socketIO := req.Header.Get("X-socket-io")
	return strings.ToLower(socketIO) == "polling"
}

func isRequestUnderTLS(req *http.Request) bool {
	return req.TLS != nil
}

func getHostname(req *http.Request) string {
	return strings.Split(req.Host, ":")[0]
}

func getServerConfig(conf config.Profile, hostname string) (config.Server, error) {
	for _, serve := range conf.Servers {
		if containString(serve.Hostname, hostname) {
			return serve, nil
		}
	}
	return config.Server{}, fmt.Errorf("hostname: %s 没有对应的转发规则", hostname)
}

func getHostnameFromTLSHelloInfo(chi *tls.ClientHelloInfo) string {
	if chi.ServerName != "" {
		return chi.ServerName
	}
	return strings.Split(chi.Conn.LocalAddr().String(), ":")[0]
}

func containString(slice []string, elem string) bool {
	for _, p := range slice {
		if p == elem {
			return true
		}
	}
	return false
}
