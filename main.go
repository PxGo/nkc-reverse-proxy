package main

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"nkc-proxy/modules"
)

func main() {
	serversPort, err := modules.GetServersPortFromConfigs()
	if err != nil {
		fmt.Println(err)
		return
	}
	ports := []int64{}
	for port, serverPort := range serversPort {
		ports = append(ports, port)
		go func(sp *modules.ServerPort) {
			_, err := modules.CreateServerAndStart(sp.Port, sp.TLSConfig)
			if err != nil {
				modules.ErrorLogger.Println(err)
				return
			}
		}(serverPort)
	}
	go func() {
		// pprof 调试
		log.Fatal(http.ListenAndServe(":6060", nil))
	}()
	fmt.Printf("server is running at %v\n", ports)
	select {}
}
