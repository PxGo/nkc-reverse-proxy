package modules

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
)

const DebugServerPort = 8080

func InitDebugServer() {
	go func() {
		err := http.ListenAndServe(fmt.Sprintf(":%d", DebugServerPort), nil)
		if err != nil {
			log.Fatal(err)
		}
	}()
}
