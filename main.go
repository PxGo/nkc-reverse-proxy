package main

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"nkc-proxy/modules"
	"strconv"
)

func main() {
	configs, err := modules.GetConfigs()
	if err != nil {
		fmt.Println(err)
		return
	}

	if configs.PProf != 0 {
		go func() {
			// pprof 调试
			fmt.Printf("pprof: localhost:%v/debug/pprof\n", configs.PProf)
			log.Fatal(http.ListenAndServe(":"+strconv.FormatInt(configs.PProf, 10), nil))
		}()
	}

	serversPort, err := modules.GetServersPortFromConfigs()
	if err != nil {
		log.Fatal(err)
	}

	httpsReverseProxy, err := modules.GetReverseProxy(true)
	if err != nil {
		log.Fatal(err)
	}
	httpReverseProxy, err := modules.GetReverseProxy(false)
	if err != nil {
		log.Fatal(err)
	}
	ports := []int64{}
	for port, serverPort := range serversPort {
		ports = append(ports, port)
		go func(sp *modules.ServerPort) {
			reverseProxy := httpReverseProxy
			if sp.TLSConfig != nil {
				reverseProxy = httpsReverseProxy
			}
			_, err := modules.CreateServerAndStart(reverseProxy, sp.Port, sp.TLSConfig)
			if err != nil {
				modules.ErrorLogger.Println(err)
				return
			}
		}(serverPort)
	}

	fmt.Printf("server is running at %v\n", ports)
	select {}
}
