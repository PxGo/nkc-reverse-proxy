package main

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"nkc-reverse-proxy/modules"
)

func main() {
	configs, err := modules.GetConfigs()
	if err != nil {
		modules.AddErrorLog(err)
		log.Fatal(err)
	}

	debugServerPort := "9527"

	if configs.Debug {
		go func() {
			fmt.Printf("Debug server is running at %v\n", debugServerPort)
			err := http.ListenAndServe(":"+debugServerPort, nil)
			if err != nil {
				modules.AddErrorLog(err)
				log.Fatal(err)
			}
		}()
	}

	serversPort, err := modules.GetServersPortFromConfigs()
	if err != nil {
		modules.AddErrorLog(err)
		log.Fatal(err)
	}

	httpsReverseProxy, err := modules.GetReverseProxy(true)
	if err != nil {
		modules.AddErrorLog(err)
		log.Fatal(err)
	}
	httpReverseProxy, err := modules.GetReverseProxy(false)
	if err != nil {
		modules.AddErrorLog(err)
		log.Fatal(err)
	}
	ports := []uint16{}
	for port, serverPort := range serversPort {
		ports = append(ports, port)
		go func(sp *modules.ServerPort) {
			reverseProxy := httpReverseProxy
			if sp.TLSConfig != nil {
				reverseProxy = httpsReverseProxy
			}
			_, err := modules.CreateServerAndStart(reverseProxy, sp.Port, sp.TLSConfig)
			if err != nil {
				modules.AddErrorLog(err)
				log.Fatal(err)
			}
		}(serverPort)
	}

	fmt.Printf("Proxy server is running at %v\n", ports)
	select {}
}
