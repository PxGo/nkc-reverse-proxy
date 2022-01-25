package modules

import (
	"log"
	"os"
)

var (
	ErrorLogger *log.Logger
)

func init() {
	errorLogPath, err := GetErrorLogPath()
	if err != nil {
		log.Fatal(err)
	}
	errorFile, err := os.OpenFile(errorLogPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	ErrorLogger = log.New(errorFile, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}
