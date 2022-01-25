package modules

import (
	"net/http"
	"net/http/httputil"
)

func GetReverseProxy(isHttps bool) (*httputil.ReverseProxy, error) {
	if isHttps && httpsReverseProxy != nil {
		return httpsReverseProxy, nil
	}
	if !isHttps && httpReverseProxy != nil {
		return httpReverseProxy, nil
	}

	director := func(req *http.Request) {
		passUrl, _, err := GetTargetPassInfo(req, isHttps)

		AddAccessLog(req.Host, req.URL, "Reverse Proxy", passUrl)

		if err != nil {
			AddErrorLog(err)
			return
		}

		req.URL.Scheme = passUrl.Scheme
		req.URL.Host = passUrl.Host

		host, _, err := GetHostInfo(req.Host, isHttps)
		if err != nil {
			AddErrorLog(err)
			return
		}

		req.Host = host
	}

	errorHandle := func(w http.ResponseWriter, r *http.Request, err error) {
		AddErrorLog(err)
		pageContent, err := GetPageByStatus(http.StatusServiceUnavailable)
		if err != nil {
			AddErrorLog(err)
			return
		}
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

	if isHttps {
		httpsReverseProxy = reverseProxy
	} else {
		httpReverseProxy = reverseProxy
	}
	return reverseProxy, nil
}
