package modules

import (
	"log"
	"os"
	"runtime/debug"
)

var (
	ErrorFileLogger  *log.Logger
	AccessFileLogger *log.Logger
	ErrorLogger      *log.Logger
	AccessLogger     *log.Logger
	Debug            bool
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

	fileLogFormat := log.Ldate | log.Ltime
	ErrorFileLogger = log.New(errorFile, "", fileLogFormat)
	ErrorLogger = log.New(os.Stderr, "ERROR ", fileLogFormat)
	AccessFileLogger = log.New(accessFile, "", fileLogFormat)
	AccessLogger = log.New(os.Stdout, "INFO ", fileLogFormat)

	Debug = configs.Debug

}

func AddAccessLog(v ...interface{}) {
	AccessFileLogger.Println(v...)
	if Debug {
		AccessLogger.Println(v...)
	}
}

func AddErrorLog(err error) {
	stackInfo := string(debug.Stack())
	ErrorFileLogger.Println(err, "\n", stackInfo)
	if Debug {
		ErrorLogger.Println(err, "\n", stackInfo)
	}
}

func AddRedirectLog(code int, url string, targetUrl string) {
	AddAccessLog("Redirect", code, url, ">>>", targetUrl)
}

func AddReverseProxyLog(url string, targetUrl string) {
	AddAccessLog("ReverseProxy", url, ">>>", targetUrl)
}

func AddNotFoundError(url string) {
	AddAccessLog("NotFound", url)
}
