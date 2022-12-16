package modules

import (
	"crypto/tls"
	"errors"
	"github.com/dgryski/go-farm"
	"github.com/go-yaml/yaml"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
)

var serverLocation ServerLocation

var configs *Configs

func GetConfigsPath() (string, string, error) {
	filePath := "configs.yaml"
	root, err := os.Getwd()
	if err != nil {
		return "", "", err
	}
	if len(os.Args) > 1 {
		filePath = os.Args[1]
	}
	if !path.IsAbs(filePath) {
		filePath = path.Join(root, filePath)
	}
	templateFilePath := path.Join(root, "configs.template.yaml")
	return filePath, templateFilePath, nil
}

func GetLogDirPath() (string, error) {
	root, err := os.Getwd()
	if err != nil {
		return "", err
	}
	logDir := path.Join(root, "./logs")
	return logDir, nil
}

func GetLogPathByLogType(logType string) (string, error) {
	logDir, err := GetLogDirPath()
	if err != nil {
		return "", err
	}
	errorLogPath := path.Join(logDir, logType+".log")
	return errorLogPath, nil
}

func GetConfigs() (*Configs, error) {
	if configs != nil {
		return configs, nil
	}
	configFilePath, templateConfigFilePath, err := GetConfigsPath()
	if err != nil {
		return nil, err
	}
	file, err := os.ReadFile(configFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			templateFile, err := ioutil.ReadFile(templateConfigFilePath)
			if err != nil {
				return nil, err
			}
			err = ioutil.WriteFile(configFilePath, templateFile, 0644)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
		file, err = os.ReadFile(configFilePath)
		if err != nil {
			return nil, err
		}
	}
	err = yaml.Unmarshal(file, &configs)
	if err != nil {
		return nil, err
	}
	return configs, nil
}

func GetServersPortFromConfigs() (map[uint16]*ServerPort, error) {
	configs, err := GetConfigs()
	if err != nil {
		return nil, err
	}
	serverPortMap := make(map[uint16]*ServerPort)
	var tlsConfig *tls.Config
	for _, server := range configs.Servers {
		var tlsCFG *tls.Config
		if server.SSLKey != "" && server.SSLCert != "" {
			if tlsConfig == nil {
				tlsConfig, err = GetTLSConfig()
				if err != nil {
					return nil, err
				}
			}
			tlsCFG = tlsConfig
		}
		serverPort := serverPortMap[server.Listen]
		if serverPort == nil {
			serverPortMap[server.Listen] = &ServerPort{Port: server.Listen, TLSConfig: tlsCFG}
		} else if serverPort.TLSConfig != tlsCFG {
			return nil, errors.New("端口冲突：" + strconv.Itoa(int(server.Listen)))
		}
	}
	return serverPortMap, nil
}

func GetTLSConfig() (*tls.Config, error) {
	cfg := tls.Config{}
	configs, err := GetConfigs()
	if err != nil {
		return nil, err
	}
	for _, server := range configs.Servers {
		if server.SSLKey == "" || server.SSLCert == "" {
			continue
		}
		cert, err := tls.LoadX509KeyPair(server.SSLCert, server.SSLKey)
		if err != nil {
			return nil, err
		}
		cfg.Certificates = append(cfg.Certificates, cert)
	}
	return &cfg, nil
}

func GetServerLocation() (ServerLocation, error) {
	if serverLocation != nil {
		return serverLocation, nil
	}
	serverLocation = make(ServerLocation)
	configs, err := GetConfigs()
	if err != nil {
		return nil, err
	}
	for _, server := range configs.Servers {
		if serverLocation[server.Listen] == nil {
			serverLocation[server.Listen] = make(NameLocation)
		}
		for _, name := range server.Name {
			if serverLocation[server.Listen][name] == nil {
				serverLocation[server.Listen][name] = server.Location
			} else {
				return nil, errors.New("端口或域名重复")
			}
		}
	}
	return serverLocation, nil
}

func GetClientIP(r *http.Request) string {
	xForwardedFor := r.Header.Get("X-Forwarded-For")
	ip := strings.TrimSpace(strings.Split(xForwardedFor, ",")[0])
	if ip != "" {
		return ip
	}

	ip = strings.TrimSpace(r.Header.Get("X-Real-Ip"))
	if ip != "" {
		return ip
	}

	if ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr)); err == nil {
		return ip
	}

	return ""
}

func GetUrlByPassType(pass []string, passType string, ip string) string {
	var index uint64
	passCount := len(pass)
	if passCount == 1 {
		index = 0
	} else if passType == "ip_hash" {
		bytes := []byte(ip)
		hash := farm.Hash64(bytes)
		index = hash % uint64(passCount)
	} else {
		index = uint64(rand.Intn(passCount))
	}
	return pass[index]
}

func GetRequestAddr(host string) (string, error) {
	hostInfo := strings.Split(host, ":")
	if len(hostInfo) == 0 {
		return "", errors.New("host error. host=" + host)
	}
	host = hostInfo[0]

	return host, nil
}

func GetTargetLocation(host string, port uint16, url string) (*Location, error) {
	serverLocation, err := GetServerLocation()
	if err != nil {
		return nil, err
	}
	if serverLocation[port] == nil ||
		serverLocation[port][host] == nil ||
		len(serverLocation[port][host]) == 0 {
		return nil, nil
	}
	locations := serverLocation[port][host]
	var targetLocation *Location
	for i := len(locations) - 1; i >= 0; i-- {
		location := locations[i]
		regString := location.Reg
		matched, err := regexp.MatchString(regString, url)
		if err != nil {
			return nil, err
		}
		if matched {
			targetLocation = &location
			break
		}
	}
	return targetLocation, nil
}

func GetLogFileByLogType(logType string) (*os.File, error) {
	fileLogPath, err := GetLogPathByLogType(logType)
	if err != nil {
		return nil, err
	}
	file, err := os.OpenFile(fileLogPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func GetLoggerByLogType(logType string) (*log.Logger, *log.Logger, error) {
	file, err := GetLogFileByLogType(logType)
	if err != nil {
		return nil, nil, err
	}
	fileLogFormat := log.Ldate | log.Ltime
	fileLogger := log.New(file, "", fileLogFormat)
	logger := log.New(os.Stderr, "["+logType+"] ", fileLogFormat)
	return fileLogger, logger, nil
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func CreateDir(path string) error {
	err := os.Mkdir(path, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func CheckAndCreateDir(path string) error {
	exist, _ := PathExists(path)
	if exist {
		return nil
	} else {
		err := CreateDir(path)
		if err != nil {
			return err
		}
		return nil
	}
}

func InitLogDir() {
	logDirPath, err := GetLogDirPath()
	if err != nil {
		log.Fatal(err)
	}
	createDirError := CheckAndCreateDir(logDirPath)
	if createDirError != nil {
		log.Fatal(createDirError)
	}
}
