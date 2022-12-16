package main

import (
	"fmt"
	"log"
	_ "net/http/pprof"
	"nkc-reverse-proxy/modules"
)

func main() {
	serversPort, err := modules.GetServersPortFromConfigs()
	if err != nil {
		log.Fatal(err)
	}
	var ports []uint16
	for port, serverPort := range serversPort {
		ports = append(ports, port)
		go func(sp *modules.ServerPort) {
			reverseProxy, err := modules.GetReverseProxy(sp.Port)
			if err != nil {
				log.Fatal(err)
			}
			_, err = modules.CreateServerAndStart(reverseProxy, sp.Port, sp.TLSConfig)
			if err != nil {
				log.Fatal(err)
			}
		}(serverPort)
	}
	fmt.Printf("Proxy server is running at %v\n", ports)
	select {}
}
