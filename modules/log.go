package modules

import (
	"log"
	"os"
)

var (
	ErrorLogger *log.Logger
)

func init() {
	configs, err := GetConfigs()
	if err != nil {
		log.Fatal(err)
	}
	errorFile, err := os.OpenFile(configs.ErrorLog, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	ErrorLogger = log.New(errorFile, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}
