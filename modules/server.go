package modules

import (
	"crypto/tls"
	"log"
	"net/http"
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

func CreateServerAndStart(port int64, cfg *tls.Config) (*http.Server, error) {
	portString := ":" + strconv.FormatInt(port, 10)
	isHttps := false
	if cfg != nil {
		isHttps = true
	}
	reverseProxy, err := GetReverseProxy(isHttps)
	if err != nil {
		return nil, err
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
		log.Fatal(server.ListenAndServeTLS("", ""))
	} else {
		log.Fatal(server.ListenAndServe())
	}
	return &server, nil
}
