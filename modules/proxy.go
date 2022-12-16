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
		location, err := GetTargetLocation(host, port, originUrl)
		if err != nil {
			AddErrorLog(err)
			return
		}
		ip := GetClientIP(req)
		targetUrlString := GetUrlByPassType(location.Pass, location.Balance, ip)
		targetUrl, err := url.Parse(targetUrlString)
		if err != nil {
			AddErrorLog(err)
			return
		}
		req.URL.Scheme = targetUrl.Scheme
		req.URL.Host = targetUrl.Host

		req.Host = originHost

		AddReverseProxyLog(req.Method, host+":"+strconv.Itoa(int(port))+originUrl, targetUrlString+originUrl)
	}

	errorHandle := func(w http.ResponseWriter, r *http.Request, err error) {
		AddErrorLog(err)
		pageContent, err := GetPageByStatus(http.StatusServiceUnavailable)
		if err != nil {
			AddErrorLog(err)
			return
		}
		w.WriteHeader(http.StatusServiceUnavailable)
		_, err = w.Write(pageContent)
		if err != nil {
			AddErrorLog(err)
			return
		}
	}

	reverseProxy := &httputil.ReverseProxy{
		Director:     director,
		ErrorHandler: errorHandle,
	}
	return reverseProxy, nil
}
