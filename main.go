package main

import (
	"fmt"
	"os"

	"github.com/tokisakiyuu/nkc-proxy-go-pure/pkg/config"
	"github.com/tokisakiyuu/nkc-proxy-go-pure/pkg/proxy"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Error: 请在命令行参数中加入配置文件路径")
	}
	configFile := os.Args[1]
	conf := config.Parse(configFile)
	serve := proxy.NewNKCProxy(conf)
	serve.Launch()
}
