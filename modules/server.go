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
		ErrorLogger.Println(err)
		return
	}
	if redirectInfo != nil && redirectInfo.Url != "" {
		http.Redirect(writer, request, redirectInfo.Url, redirectInfo.Code)
	} else if passUrl != nil {
		handle.ReverseProxy.ServeHTTP(writer, request)
	} else {
		pageContent, err := GetPageByStatus(http.StatusNotFound)
		if err != nil {
			ErrorLogger.Println(err)
			return
		}
		writer.WriteHeader(http.StatusNotFound)
		_, err = writer.Write(pageContent)
		if err != nil {
			ErrorLogger.Println(err)
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
			ErrorLogger.Println(err)
			log.Fatal(err)
		}
	} else {
		err := server.ListenAndServe()
		if err != nil {
			ErrorLogger.Println(err)
			log.Fatal(err)
		}
	}
	return &server, nil
}
