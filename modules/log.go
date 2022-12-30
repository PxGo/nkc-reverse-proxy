package modules

import (
	"log"
	"os"
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

	ErrorFileLogger, ErrorLogger, err = GetLoggerByLogType("error", os.Stderr)
	if err != nil {
		log.Fatal(err)
	}
	InfoFileLogger, InfoLogger, err = GetLoggerByLogType("info", os.Stdout)
	if err != nil {
		log.Fatal(err)
	}
	WarningFileLogger, WarningLogger, err = GetLoggerByLogType("warning", os.Stderr)
	if err != nil {
		log.Fatal(err)
	}
	DebugFileLogger, DebugLogger, err = GetLoggerByLogType("debug", os.Stdout)
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

func AddRedirectLog(ip string, port string, method string, code int, url string, targetUrl string) {
	AddInfoLog("[" + ip + ":" + port + "] " + "Redirect" + " " + method + " " + string(rune(code)) + " " + url + " " + targetUrl)
}

func AddReverseProxyLog(ip string, port string, method string, url string, targetUrl string) {
	AddInfoLog("[" + ip + ":" + port + "] " + "ReverseProxy" + " " + method + " " + url + " " + ">>>" + " " + targetUrl)
}

func AddNotFoundError(ip string, port string, method string, url string) {
	AddInfoLog("[" + ip + ":" + port + "] " + "NotFound" + " " + method + " " + url)
}
