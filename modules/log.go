package modules

import (
	"fmt"
	"log"
	"os"
)

var (
	ErrorLogger  *log.Logger
	AccessLogger *log.Logger
	Debug        bool
)

func init() {

	configs, err := GetConfigs()
	if err != nil {
		log.Fatal(err)
	}

	errorLogPath, err := GetLogPathByLogType("error")
	if err != nil {
		log.Fatal(err)
	}
	accessLogPath, err := GetLogPathByLogType("access")
	if err != nil {
		log.Fatal(err)
	}
	errorFile, err := os.OpenFile(errorLogPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	accessFile, err := os.OpenFile(accessLogPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	ErrorLogger = log.New(errorFile, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	AccessLogger = log.New(accessFile, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)

	Debug = configs.Debug

}

func AddAccessLog(v ...interface{}) {
	AccessLogger.Println(v)
	if Debug {
		fmt.Println(v)
	}
}

func AddErrorLog(v ...interface{}) {
	ErrorLogger.Println(v)
	if Debug {
		fmt.Println(v)
	}
}
