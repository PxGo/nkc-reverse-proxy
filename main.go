package main

import (
	"fmt"
	"log"
	_ "net/http/pprof"
	"nkc-reverse-proxy/modules"
)

func main() {

	err := modules.InitGlobalConfigs()
	if err != nil {
		modules.AddErrorLog(err)
		log.Fatal(err)
	}
	err = modules.InitGlobalServices()
	if err != nil {
		modules.AddErrorLog(err)
		log.Fatal(err)
	}
	err = modules.InitAutoCert()
	if err != nil {
		modules.AddErrorLog(err)
		log.Fatal(err)
	}

	serversPort, err := modules.GetServersPortFromConfigs()
	if err != nil {
		modules.AddErrorLog(err)
		log.Fatal(err)
	}
	var ports []uint16
	for port, serverPort := range serversPort {
		ports = append(ports, port)
		go func(sp *modules.ServerPort) {
			reverseProxy, err := modules.GetReverseProxy(sp.Port)
			if err != nil {
				modules.AddErrorLog(err)
				log.Fatal(err)
			}
			_, err = modules.CreateServerAndStart(reverseProxy, sp.Port, sp.TLSConfig)
			if err != nil {
				modules.AddErrorLog(err)
				log.Fatal(err)
			}
		}(serverPort)
	}
	fmt.Printf("Proxy server is running at %v\n", ports)

	if modules.GlobalConfigs.Console.Debug {
		modules.InitDebugServer()
		fmt.Printf("Debug server is running at %d\n", modules.DebugServerPort)
	}

	select {}
}
