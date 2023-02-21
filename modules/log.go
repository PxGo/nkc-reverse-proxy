package modules

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"time"
)

type ILoggerType string
type ILoggerContent struct {
	fileLogger    *log.Logger
	consoleLogger *log.Logger
	file          *os.File
}
type ILoggers map[ILoggerType]*ILoggerContent

type ILoggerChanData struct {
	logType ILoggerType
	content string
}
type ILoggerChan chan ILoggerChanData

type IConsole map[ILoggerType]bool

const (
	LoggerTypeError ILoggerType = "error"
	LoggerTypeInfo  ILoggerType = "info"
	LoggerTypeDebug ILoggerType = "debug"
	LoggerTypeWarn  ILoggerType = "warn"
)

type Logger struct {
	date    string
	loggers ILoggers
	logChan ILoggerChan
	console IConsole
}

func (logger *Logger) createLogger(logType ILoggerType, date string, loggerStd io.Writer) (*log.Logger, *log.Logger, *os.File, error) {
	logFilePath, err := GetLogPathByLogType(logType, date)
	if err != nil {
		return nil, nil, nil, err
	}
	file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return nil, nil, nil, err
	}
	fileLogFormat := log.Ldate | log.Ltime
	fileLogger := log.New(file, "", fileLogFormat)
	consoleLogger := log.New(loggerStd, "["+string(logType)+"] ", fileLogFormat)
	return consoleLogger, fileLogger, file, nil
}

func (logger *Logger) initLoggers(date string) error {
	loggers := make(ILoggers)
	logTypes := []ILoggerType{LoggerTypeError, LoggerTypeDebug, LoggerTypeWarn, LoggerTypeInfo}
	for _, logType := range logTypes {
		var console io.Writer = os.Stdout
		if logType == LoggerTypeError {
			console = os.Stderr
		}
		consoleLogger, fileLogger, file, err := logger.createLogger(logType, date, console)
		if err != nil {
			return err
		}
		if logger.loggers[logType] != nil {
			err := logger.loggers[logType].file.Close()
			if err != nil {
				fmt.Println(err)
			}
		}

		loggers[logType] = &ILoggerContent{
			consoleLogger: consoleLogger,
			fileLogger:    fileLogger,
			file:          file,
		}
	}
	logger.date = date
	logger.loggers = loggers
	return nil
}

func (logger *Logger) getLoggerByType(logType ILoggerType) (*ILoggerContent, error) {
	date := time.Now().Format("2006-01-02")
	if logger.date != date {
		err := logger.initLoggers(date)
		if err != nil {
			return nil, err
		}
	}
	return logger.loggers[logType], nil
}

func (logger *Logger) Println(logType ILoggerType, content string) {
	loggerContent, err := logger.getLoggerByType(logType)
	if err != nil {
		fmt.Println(content)
		fmt.Println(err)
		return
	}
	consoleConfig := logger.console
	if consoleConfig[logType] {
		loggerContent.consoleLogger.Println(content)
	}
	loggerContent.fileLogger.Println(content)
}

func (logger *Logger) SendLogData(logType ILoggerType, content string) {
	go func(logType ILoggerType, content string) {
		logger.logChan <- ILoggerChanData{
			logType: logType,
			content: content,
		}
	}(logType, content)
}

func (logger *Logger) InfoLog(content string) {
	logger.SendLogData(LoggerTypeInfo, content)
}

func (logger *Logger) ErrorLog(content string) {
	logger.SendLogData(LoggerTypeError, content)
}

func (logger *Logger) WarnLog(content string) {
	logger.SendLogData(LoggerTypeWarn, content)
}

func (logger *Logger) DebugLog(content string) {
	logger.SendLogData(LoggerTypeDebug, content)
}

var logger Logger

func init() {

	configs, err := GetConfigs()
	if err != nil {
		log.Fatal(err)
	}

	InitLogDir()
	logger = Logger{
		date:    "",
		loggers: make(ILoggers),
		logChan: make(ILoggerChan),
		console: IConsole{
			LoggerTypeWarn:  configs.Console.Warning,
			LoggerTypeError: configs.Console.Error,
			LoggerTypeDebug: configs.Console.Debug,
			LoggerTypeInfo:  configs.Console.Info,
		},
	}
	go func(logger *Logger) {
		for {
			data := <-logger.logChan
			logger.Println(data.logType, data.content)
		}
	}(&logger)
}

func GetLogDirPath() (string, error) {
	root, err := os.Getwd()
	if err != nil {
		return "", err
	}
	logDir := path.Join(root, "./logs")
	return logDir, nil
}

func GetLogPathByLogType(logType ILoggerType, date string) (string, error) {
	logDir, err := GetLogDirPath()
	if err != nil {
		return "", err
	}
	errorLogPath := path.Join(logDir, date+"."+string(logType)+".log")
	return errorLogPath, nil
}

func AddErrorLog(err error) {
	logger.ErrorLog(err.Error() /*"\n", stackInfo*/)
}

func AddInfoLog(content string) {
	logger.InfoLog(content)
}

/*func AddWarningLog(content string) {
	logger.WarnLog(content)
}

func AddDebugLog(content string) {
	logger.DebugLog(content)
}*/

func AddRedirectLog(ip string, port string, method string, code int, url string, targetUrl string) {
	content := fmt.Sprintf("[%s:%s] Redirect %s %s %s %s", ip, port, method, string(rune(code)), url, targetUrl)
	AddInfoLog(content)
}

func AddReverseProxyLog(ip string, port string, method string, url string, targetUrl string) {
	content := fmt.Sprintf("[%s:%s] ReverseProxy %s %s >>> %s", ip, port, method, url, targetUrl)
	AddInfoLog(content)
}

func AddNotFoundError(ip string, port string, method string, url string) {
	content := fmt.Sprintf("[%s:%s] NotFound %s %s", ip, port, method, url)
	AddInfoLog(content)
}

func AddServiceUnavailableError(ip string, port string, method string, url string) {
	content := fmt.Sprintf("[%s:%s] ServiceUnavailable %s %s", ip, port, method, url)
	AddInfoLog(content)
}

func AddReqLimitInfo(ip string, port string, method string, url string, reqLimitType string) {
	content := fmt.Sprintf("[%s:%s] TooManyRequest %s %s %s", ip, port, reqLimitType, method, url)
	AddInfoLog(content)
}
