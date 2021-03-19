package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/tokisakiyuu/nkc-proxy-go-pure/pkg/balance"
	"github.com/tokisakiyuu/nkc-proxy-go-pure/pkg/config"
)

var ReverseProxyCaches = make(map[string]*httputil.ReverseProxy)

func responseWebSocket(conf config.Server, rw http.ResponseWriter, req *http.Request) {
	balancer := balance.NewIPKeyBalancer()
	balancer.Add(conf.WsTarget...)
	hostname := getHostname(req)
	targetHost := balancer.Get(hostname)
	// 转发
	var proxy *httputil.ReverseProxy
	if ReverseProxyCaches[targetHost] != nil {
		proxy = ReverseProxyCaches[targetHost]
	} else {
		remote, _ := url.Parse(targetHost)
		proxy = httputil.NewSingleHostReverseProxy(remote)
		ReverseProxyCaches[targetHost] = proxy
	}
	proxy.ServeHTTP(rw, req)
}

func responseHTTP(conf config.Server, rw http.ResponseWriter, req *http.Request) {
	balancer := balance.NewRandomBalancer()
	balancer.Add(conf.HttpTarget...)
	targetHost := balancer.Get()
	// 如果是重定向
	if conf.Type == "redirect" {
		http.Redirect(rw, req, targetHost+req.RequestURI, http.StatusMovedPermanently)
		return
	}
	// 否则就是转发
	remote, _ := url.Parse(targetHost)
	var proxy *httputil.ReverseProxy
	if ReverseProxyCaches[targetHost] != nil {
		proxy = ReverseProxyCaches[targetHost]
		proxy.ErrorHandler = func(rw http.ResponseWriter, req *http.Request, e error) {
			responseError(conf, rw, req)
		}
	} else {
		proxy = httputil.NewSingleHostReverseProxy(remote)
		ReverseProxyCaches[targetHost] = proxy
	}
	req.URL.Scheme = remote.Scheme
	req.URL.Host = remote.Host
	req.Host = remote.Host
	proxy.ServeHTTP(rw, req)
}

// 响应一个代理错误页面，并非来自后端服务器响应的错误
func responseError(conf config.Server, rw http.ResponseWriter, req *http.Request) {
	if conf.NoResponsePage != "" {
		http.ServeFile(rw, req, conf.NoResponsePage)
		return
	}
	rw.WriteHeader(403)
	rw.Write([]byte("Forbidden 服务器无响应"))
}
