package modules

import (
	"log"
)

var (
	ErrorFileLogger   *log.Logger
	ErrorLogger       *log.Logger
	InfoFileLogger    *log.Logger
	InfoLogger        *log.Logger
	WarningFileLogger *log.Logger
	WarningLogger     *log.Logger
	DebugFileLogger   *log.Logger
	DebugLogger       *log.Logger
)

var console Console

func init() {

	configs, err := GetConfigs()
	if err != nil {
		log.Fatal(err)
	}

	InitLogDir()

	console.Debug = configs.Console.Debug
	console.Info = configs.Console.Info
	console.Warning = configs.Console.Warning
	console.Error = configs.Console.Error

	ErrorFileLogger, ErrorLogger, err = GetLoggerByLogType("error")
	if err != nil {
		log.Fatal(err)
	}
	InfoFileLogger, InfoLogger, err = GetLoggerByLogType("info")
	if err != nil {
		log.Fatal(err)
	}
	WarningFileLogger, WarningLogger, err = GetLoggerByLogType("warning")
	if err != nil {
		log.Fatal(err)
	}
	DebugFileLogger, DebugLogger, err = GetLoggerByLogType("debug")
	if err != nil {
		log.Fatal(err)
	}
}

func AddErrorLog(err error) {
	ErrorFileLogger.Println(err /*"\n", stackInfo*/)
	if console.Error {
		ErrorLogger.Println(err /*"\n", stackInfo*/)
	}
}

func AddInfoLog(content string) {
	InfoFileLogger.Println(content)
	if console.Info {
		InfoLogger.Println(content)
	}
}

func AddWarningLog(content string) {
	WarningFileLogger.Println(content)
	if console.Warning {
		WarningLogger.Println(content)
	}
}

func AddDebugLog(content string) {
	DebugFileLogger.Println(content)
	if console.Debug {
		DebugLogger.Println(content)
	}
}

func AddRedirectLog(method string, code int, url string, targetUrl string) {
	AddInfoLog("Redirect" + " " + method + " " + string(rune(code)) + " " + url + " " + targetUrl)
}

func AddReverseProxyLog(method string, url string, targetUrl string) {
	AddInfoLog("ReverseProxy" + " " + method + " " + url + " " + ">>>" + " " + targetUrl)
}

func AddNotFoundError(method string, url string) {
	AddInfoLog("NotFound" + " " + method + " " + url)
}
