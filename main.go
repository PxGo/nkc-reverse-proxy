package main

import (
	"fmt"
	"nkc-proxy/tools"
)

func main() {
	serversPort, err := tools.GetServersPortFromConfigs()
	if err != nil {
		fmt.Println(err)
		return
	}
	messages := make(chan string)
	ports := []int64{}
	for port, serverPort := range serversPort {
		ports = append(ports, port)
		go func(sp *tools.ServerPort) {
			_, err := tools.CreateServerAndStart(sp.Port, sp.TLSConfig)
			if err != nil {
				fmt.Printf("创建服务出错")
				fmt.Printf(err.Error())
				return
			}
		}(serverPort)
	}
	fmt.Printf("server is running at %v", ports)
	<-messages
}
