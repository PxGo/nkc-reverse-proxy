/*package main

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
	fmt.Println(conf)
	serve := proxy.NewNKCProxy(conf)
	serve.Launch()
}
*/

package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
)

func homeHandle(w http.ResponseWriter, r *http.Request) {
	_, err := io.WriteString(w, "hello, world")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

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

func main() {
	proxy := newReverseProxy()
	go func() {
		log.Fatal(http.ListenAndServe(":8080", proxy))
	}()
	log.Fatal(http.ListenAndServe(":9090", proxy))
}
