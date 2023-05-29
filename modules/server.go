package modules

import (
	"crypto/tls"
	"errors"
	"net/http"
	"net/http/httputil"
	"net/url"
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
		err := WriteResponse(request, writer, http.StatusInternalServerError, GlobalConfigs.Template.Page500)
		if err != nil {
			AddErrorLog(err)
		}
		return
	}

	service, err := GetTargetService(host, handle.Port, request.URL.String())
	if err != nil {
		AddErrorLog(err)
		err := WriteResponse(request, writer, http.StatusInternalServerError, GlobalConfigs.Template.Page500)
		if err != nil {
			AddErrorLog(err)
		}
		return
	}

	if service == nil {
		// 不存在匹配的服务
		// 返回404
		AddNotFoundError(ip, port, request.Method, request.Host+request.URL.String())
		err := WriteResponse(request, writer, http.StatusNotFound, GlobalConfigs.Template.Page404)
		if err != nil {
			AddErrorLog(err)
		}
		return
	}

	// 存在匹配的服务

	// 需要重定向
	if service.Location.RedirectUrl != "" && service.Location.RedirectCode != 0 {
		redirectUrl, err := url.Parse(service.Location.RedirectUrl)
		if err != nil {
			AddErrorLog(errors.New("redirect url parse error. url=" + service.Location.RedirectUrl))
			return
		}
		if redirectUrl.Path == "" {
			redirectUrl.Path = request.URL.Path
		}
		AddRedirectLog(ip, port, request.Method, service.Location.RedirectCode, request.Host+request.URL.String(), redirectUrl.String())
		http.Redirect(writer, request, redirectUrl.String(), service.Location.RedirectCode)
		return
	}

	if service.Location.Pass == nil || len(service.Location.Pass) == 0 {
		// 目标服务为空
		AddServiceUnavailableError(ip, port, request.Method, request.Host+request.URL.String())
		err := WriteResponse(request, writer, http.StatusServiceUnavailable, service.Template.Page503)
		if err != nil {
			AddErrorLog(err)
		}
		return
	}

	limited := ReqLimitChecker(service.Global.ReqLimit, ip)
	if limited {
		AddReqLimitInfo(ip, port, request.Method, request.URL.String(), "Global")
		err := WriteResponse(request, writer, http.StatusTooManyRequests, service.Template.Page429)
		if err != nil {
			AddErrorLog(err)
		}
		return
	}

	limited = ReqLimitChecker(service.Server.ReqLimit, ip)
	if limited {
		AddReqLimitInfo(ip, port, request.Method, request.URL.String(), "Server")
		err := WriteResponse(request, writer, http.StatusTooManyRequests, service.Template.Page429)
		if err != nil {
			AddErrorLog(err)
		}
		return
	}

	limited = ReqLimitChecker(service.Location.ReqLimit, ip)
	if limited {
		AddReqLimitInfo(ip, port, request.Method, request.URL.String(), "Location")
		err := WriteResponse(request, writer, http.StatusTooManyRequests, service.Template.Page429)
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

	var targetHandle http.Handler = &NKCHandle{
		IsHTTPS:      isHttps,
		Port:         port,
		ReverseProxy: reverseProxy,
	}

	if !isHttps && AutoCertIsEnabled() {
		targetHandle = AutoCert.HTTPHandler(targetHandle)
	}

	server := http.Server{
		Addr:    portString,
		Handler: targetHandle,
	}
	if isHttps {
		server.TLSConfig = cfg
		err := server.ListenAndServeTLS("", "")
		if err != nil {
			return nil, err
		}
	} else {
		err := server.ListenAndServe()
		if err != nil {
			return nil, err
		}
	}
	return &server, nil
}
