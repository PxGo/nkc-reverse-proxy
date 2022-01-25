package modules

import (
	"crypto/tls"
	"log"
	"net/http"
	"net/http/httputil"
	"strconv"
)

func (handle NKCHandle) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	passUrl, redirectInfo, err := GetTargetPassInfo(request, handle.IsHTTPS)
	if err != nil {
		AddErrorLog(err)
		return
	}
	if redirectInfo != nil && redirectInfo.Url != "" {
		AddAccessLog(request.Host, request.URL, "Redirect", redirectInfo.Code, redirectInfo.Url)
		http.Redirect(writer, request, redirectInfo.Url, redirectInfo.Code)
	} else if passUrl != nil {
		handle.ReverseProxy.ServeHTTP(writer, request)
	} else {
		AddAccessLog(request.Host, request.URL, "NotFound")
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
