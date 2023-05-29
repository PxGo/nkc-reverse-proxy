package main

import (
	"fmt"
	_ "net/http/pprof"
	"nkc-reverse-proxy/modules"
	"os"
	"time"
)

func main() {

	err := modules.InitGlobalConfigs()
	if err != nil {
		terminal(err)
	}
	err = modules.InitGlobalServices()
	if err != nil {
		terminal(err)
	}
	err = modules.InitAutoCert()
	if err != nil {
		terminal(err)
	}

	serversPort, err := modules.GetServersPortFromConfigs()
	if err != nil {
		terminal(err)
	}
	var ports []uint16
	for port, serverPort := range serversPort {
		ports = append(ports, port)
		go func(sp *modules.ServerPort) {
			reverseProxy, err := modules.GetReverseProxy(sp.Port)
			if err != nil {
				terminal(err)
			}
			_, err = modules.CreateServerAndStart(reverseProxy, sp.Port, sp.TLSConfig)
			if err != nil {
				terminal(err)
			}
		}(serverPort)
	}
	fmt.Printf("NRP[%v] is running at %v\n", modules.CodeVersion, ports)

	if modules.GlobalConfigs.Console.Debug {
		modules.InitDebugServer()
		fmt.Printf("Debug server is running at %d\n", modules.DebugServerPort)
	}

	select {}
}

func terminal(err error) {
	modules.AddErrorLog(err)
	time.Sleep(time.Second)
	os.Exit(1)
}
