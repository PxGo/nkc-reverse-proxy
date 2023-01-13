package modules

import (
	"crypto/tls"
	"errors"
	"log"
	"net/http"
	"net/http/httputil"
	"strconv"
)

type ServerPort struct {
	Port      uint16
	TLSConfig *tls.Config
}

type NKCHandle struct {
	IsHTTPS      bool
	Port         uint16
	ReverseProxy *httputil.ReverseProxy
}

func (handle NKCHandle) ServeHTTP(writer http.ResponseWriter, request *http.Request) {

	ip, port := GetClientRealAddr(request)

	// 获取请求的域、端口以及路径
	host, err := GetRequestAddr(request.Host)

	if err != nil {
		AddErrorLog(err)
		err := WriteResponse(request, writer, http.StatusInternalServerError)
		if err != nil {
			AddErrorLog(err)
		}
		return
	}

	service, err := GetTargetService(host, handle.Port, request.URL.String())

	if service == nil {
		// 不存在匹配的服务
		// 返回404
		AddNotFoundError(ip, port, request.Method, request.Host+request.URL.String())
		err := WriteResponse(request, writer, http.StatusNotFound)
		if err != nil {
			AddErrorLog(err)
		}
		return
	}

	if service.Location.Pass == nil || len(service.Location.Pass) == 0 {
		// 目标服务为空
		AddServiceUnavailableError(ip, port, request.Method, request.Host+request.URL.String())
		err := WriteResponse(request, writer, http.StatusServiceUnavailable)
		if err != nil {
			AddErrorLog(err)
		}
		return
	}

	// 存在匹配的服务

	limited := ReqLimitChecker(service.Global.ReqLimit, ip)
	if limited {
		AddErrorLog(errors.New("global req limit: too Many Request"))
		err := WriteResponse(request, writer, http.StatusTooManyRequests)
		if err != nil {
			AddErrorLog(err)
		}
		return
	}

	limited = ReqLimitChecker(service.Server.ReqLimit, ip)
	if limited {
		AddErrorLog(errors.New("global req limit: too Many Request"))
		err := WriteResponse(request, writer, http.StatusTooManyRequests)
		if err != nil {
			AddErrorLog(err)
		}
		return
	}

	limited = ReqLimitChecker(service.Location.ReqLimit, ip)
	if limited {
		AddErrorLog(errors.New("global req limit: too Many Request"))
		err := WriteResponse(request, writer, http.StatusTooManyRequests)
		if err != nil {
			AddErrorLog(err)
		}
		return
	}

	handle.ReverseProxy.ServeHTTP(writer, request)
}

func CreateServerAndStart(reverseProxy *httputil.ReverseProxy, port uint16, cfg *tls.Config) (*http.Server, error) {
	portString := ":" + strconv.Itoa(int(port))
	isHttps := false
	if cfg != nil {
		isHttps = true
	}
	server := http.Server{
		Addr: portString,
		Handler: &NKCHandle{
			IsHTTPS:      isHttps,
			Port:         port,
			ReverseProxy: reverseProxy,
		},
	}
	if isHttps {
		server.TLSConfig = cfg
		err := server.ListenAndServeTLS("", "")
		if err != nil {
			AddErrorLog(err)
			log.Fatal(err)
		}
	} else {
		err := server.ListenAndServe()
		if err != nil {
			AddErrorLog(err)
			log.Fatal(err)
		}
	}
	return &server, nil
}
