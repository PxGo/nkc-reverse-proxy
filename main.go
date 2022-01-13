package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"nkc-proxy/pkg/tools"
	"strconv"
)

func newReverseProxy() *httputil.ReverseProxy {
	director := func(req *http.Request) {
		req.URL.Scheme = "http"
		req.URL.Host = "localhost:9000"
		//req.URL.Path = req
	}
	return &httputil.ReverseProxy{
		Director: director,
	}
}

func CreateServer(proxyServer *httputil.ReverseProxy, port int64, SSLKey string, SSLCert string) {
	hasSSL := SSLKey != "" && SSLCert != ""
	portString := ":" + strconv.FormatInt(port, 10)
	if hasSSL {
		log.Fatal(http.ListenAndServeTLS(portString, SSLCert, SSLKey, proxyServer))
	} else {
		log.Fatal(http.ListenAndServe(portString, proxyServer))
	}
}

func main() {
	serversPort, err := tools.GetPortFromConfigs()
	if err != nil {
		fmt.Println(err)
		return
	}
	proxy := newReverseProxy()
	go func() {
		for index, serverPort := range serversPort {
			if index == 0 {
				continue
			}
			log.Fatal(http.ListenAndServe(":"+strconv.FormatInt(serverPort, 10), proxy))
		}
	}()
	log.Fatal(http.ListenAndServe(":"+strconv.FormatInt(serversPort[0], 10), proxy))
}
