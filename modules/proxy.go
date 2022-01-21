package modules

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"time"
)

func GetReverseProxy(isHttps bool) (*httputil.ReverseProxy, error) {
	if isHttps && httpsReverseProxy != nil {
		return httpsReverseProxy, nil
	}
	if !isHttps && httpReverseProxy != nil {
		return httpReverseProxy, nil
	}

	configs, err := GetConfigs()
	if err != nil {
		return nil, err
	}
	transportConfig := configs.Transport

	fmt.Println(transportConfig, "\n")

	transport := &http.Transport{
		Proxy:             http.ProxyFromEnvironment,
		DisableKeepAlives: !transportConfig.KeepAlive,
		DialContext: (&net.Dialer{
			Timeout:   transportConfig.Timeout * time.Second,
			KeepAlive: transportConfig.KeepAliveTimeout * time.Second,
		}).DialContext,
		ForceAttemptHTTP2:   false,
		MaxIdleConns:        transportConfig.MaxIdleConnections,
		IdleConnTimeout:     transportConfig.IdleConnectionTimeout * time.Second,
		MaxIdleConnsPerHost: transportConfig.MaxIdleConnectionsPerHost,
		MaxConnsPerHost:     transportConfig.MaxConnectionsPerHost,
	}

	director := func(req *http.Request) {
		passUrl, _, err := GetTargetPassInfo(req, isHttps)
		if err != nil {
			ErrorLogger.Println(err)
			return
		}
		req.URL.Scheme = passUrl.Scheme
		req.URL.Host = passUrl.Host
	}

	errorHandle := func(w http.ResponseWriter, r *http.Request, err error) {
		ErrorLogger.Println(err)
		pageContent, err := GetPageByStatus(http.StatusServiceUnavailable)
		if err != nil {
			ErrorLogger.Println(err)
			return
		}
		_, err = w.Write(pageContent)
		if err != nil {
			ErrorLogger.Println(err)
			return
		}
	}

	reverseProxy := &httputil.ReverseProxy{
		Transport:    transport,
		Director:     director,
		ErrorHandler: errorHandle,
	}

	if isHttps {
		httpsReverseProxy = reverseProxy
	} else {
		httpReverseProxy = reverseProxy
	}
	return reverseProxy, nil
}
