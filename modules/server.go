package modules

import (
	"crypto/tls"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
)

func (handle NKCHandle) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	// 获取请求的域、端口以及路径
	host, err := GetRequestAddr(request.Host)
	location, err := GetTargetLocation(host, handle.Port, request.URL.String())
	if err != nil {
		AddErrorLog(err)
		return
	}
	if location != nil && location.RedirectUrl != "" && location.RedirectCode != 0 {
		// 重定向
		redirectUrl, err := url.Parse(location.RedirectUrl)
		if err != nil {
			AddErrorLog(err)
			return
		}
		if redirectUrl.Path == "" {
			redirectUrl.Path = request.URL.Path
		}
		ip, port := GetClientRealAddr(request)
		AddRedirectLog(ip, port, request.Method, location.RedirectCode, request.Host+request.URL.String(), redirectUrl.String())
		http.Redirect(writer, request, redirectUrl.String(), location.RedirectCode)
	} else if location != nil && location.Pass != nil && len(location.Pass) > 0 {
		handle.ReverseProxy.ServeHTTP(writer, request)
	} else {
		ip, port := GetClientRealAddr(request)
		AddNotFoundError(ip, port, request.Method, request.Host+request.URL.String())
		pageContent, err := GetPageByStatus(http.StatusNotFound)
		if err != nil {
			AddErrorLog(err)
			return
		}
		writer.WriteHeader(http.StatusNotFound)
		_, err = writer.Write(pageContent)
		if err != nil {
			AddErrorLog(err)
			return
		}
	}
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
