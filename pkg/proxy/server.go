package proxy

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"runtime/debug"
	"strconv"
	"time"

	"github.com/tokisakiyuu/nkc-proxy-go-pure/pkg/config"
)

// 创建http协议代理服务器
func NewServer(conf config.Profile) *http.Server {
	var (
		tlsConfig  = &tls.Config{}
		serveMux   = http.NewServeMux()
		concurrent = NewConcurrentLimit(conf.ConcurrentLimit)
	)
	// cert, _ := tls.LoadX509KeyPair(`D:\zlp\nkc-proxy-go-pure\assets\cert\www.kechuang.org.crt`, `D:\zlp\nkc-proxy-go-pure\assets\cert\www.kechuang.org.key`)
	tlsConfig.GetCertificate = func(chi *tls.ClientHelloInfo) (*tls.Certificate, error) {
		serveConf, _ := getServerConfig(conf, getHostnameFromTLSHelloInfo(chi))
		cert, err := config.GetCertificate(serveConf.SSL.Cert)
		if err != nil {
			fmt.Printf("%s\n%s\n", err.Error(), debug.Stack())
			return nil, err
		}
		return cert, nil
	}

	serveMux.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {
		concurrent.Add()
		defer concurrent.Pass()
		// 获取到对应的server配置
		serveConf, err := getServerConfig(conf, getHostname(req))
		if err != nil {
			fmt.Printf("%s\n%s\n", err.Error(), debug.Stack())
			rw.WriteHeader(500)
			rw.Write([]byte("代理中未配置转发规则"))
			return
		}
		// 非https请求的话，根据需要进行重定向
		if serveConf.Https && !isRequestUnderTLS(req) {
			hostname := getHostname(req)
			http.Redirect(rw, req, "https://"+hostname+":"+strconv.FormatInt(conf.Ports.HttpsPort, 10)+req.RequestURI, http.StatusMovedPermanently)
			return
		}
		// 是websocket请求的话
		if isWebSocketUpgradeRequest(req) {
			responseWebSocket(serveConf, rw, req)
			return
		}
		// 是socket.io的轮询请求的话
		if isSocketIOPolling(req) {
			responseWebSocket(serveConf, rw, req)
			return
		}
		// 是普通http请求的话
		responseHTTP(serveConf, rw, req)
	})

	return &http.Server{
		Addr:         ":" + strconv.FormatInt(conf.Ports.HttpPort, 10),
		Handler:      serveMux,
		WriteTimeout: time.Duration(conf.Timeout * int64(time.Millisecond)),
		TLSConfig:    tlsConfig,
	}
}

// golang内置反向代理handle实例的缓存
var ReverseProxyCaches = make(map[string]*httputil.ReverseProxy)

// 获取反向代理handle实例
func GetReverseProxyer(host string) *httputil.ReverseProxy {
	if ReverseProxyCaches[host] != nil {
		return ReverseProxyCaches[host]
	} else {
		remote, _ := url.Parse(host)
		proxy := httputil.NewSingleHostReverseProxy(remote)
		defer func() {
			ReverseProxyCaches[host] = proxy
		}()
		return proxy
	}
}
