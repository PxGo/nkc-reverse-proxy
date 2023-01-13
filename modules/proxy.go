package modules

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
)

func GetReverseProxy(port uint16) (*httputil.ReverseProxy, error) {
	director := func(req *http.Request) {
		originHost := req.Host
		host, err := GetRequestAddr(originHost)
		originUrl := req.URL.String()
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

		SetXForwardedRemotePort(req)

		AddReverseProxyLog(realIp, realPort, req.Method, host+":"+strconv.Itoa(int(port))+originUrl, targetUrlString+originUrl)
	}

	errorHandle := func(w http.ResponseWriter, r *http.Request, err error) {
		ip, port := GetClientRealAddr(r)
		AddErrorLog(err)
		AddServiceUnavailableError(ip, port, r.Method, r.URL.String())
		err = WriteResponse(r, w, http.StatusServiceUnavailable)
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
