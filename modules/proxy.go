package modules

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
)

const (
	HeaderKeyOriginUrl  = "X-Proxy-Origin-Url"
	HeaderKeyOriginHost = "X-Proxy-Origin-Host"
)

func GetReverseProxy(port uint16) (*httputil.ReverseProxy, error) {
	director := func(req *http.Request) {
		originHost := req.Host
		originUrl := req.URL.String()
		host, err := GetRequestAddr(originHost)
		if err != nil {
			AddErrorLog(err)
			return
		}
		service, err := GetTargetService(host, port, originUrl)
		if err != nil {
			AddErrorLog(err)
			return
		}
		location := service.Location
		realIp, realPort := GetClientRealAddr(req)
		targetUrlString := GetUrlByPassType(location.Pass, location.Balance, realIp)
		targetUrl, err := url.Parse(targetUrlString)
		if err != nil {
			AddErrorLog(err)
			return
		}
		req.URL.Scheme = targetUrl.Scheme
		req.URL.Host = targetUrl.Host

		req.Host = originHost

		req.Header.Set(HeaderKeyOriginUrl, originUrl)
		req.Header.Set(HeaderKeyOriginHost, originHost)

		SetXForwardedRemotePort(req)

		AddReverseProxyLog(realIp, realPort, req.Method, host+":"+strconv.Itoa(int(port))+originUrl, targetUrlString+originUrl)
	}

	errorHandle := func(w http.ResponseWriter, r *http.Request, err error) {
		AddErrorLog(err)

		originUrl := r.Header.Get(HeaderKeyOriginUrl)
		originHost := r.Header.Get(HeaderKeyOriginHost)

		service, err := GetTargetService(originHost, port, originUrl)

		if err != nil {
			AddErrorLog(err)
			return
		}

		ip, rPort := GetClientRealAddr(r)

		AddServiceUnavailableError(ip, rPort, r.Method, r.URL.String())

		err = WriteResponse(r, w, http.StatusServiceUnavailable, service.Template.Page503)
		if err != nil {
			AddErrorLog(err)
		}
	}

	reverseProxy := &httputil.ReverseProxy{
		Director:     director,
		ErrorHandler: errorHandle,
	}
	return reverseProxy, nil
}
