package modules

import (
	"crypto/tls"
	"errors"
	"github.com/dgryski/go-farm"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const IpHeader = "X-Forwarded-For"
const PortHeader = "X-Forwarded-Remote-Port"

func GetServersPortFromConfigs() (map[uint16]*ServerPort, error) {
	configs := GlobalConfigs
	serverPortMap := make(map[uint16]*ServerPort)

	TLSConfig, err := GetTLSConfig(false)
	if err != nil {
		return nil, err
	}
	autoTLSConfig, err := GetTLSConfig(true)
	if err != nil {
		return nil, err
	}

	for _, server := range configs.Servers {
		var tlsCFG *tls.Config

		if server.SSLAuto {
			tlsCFG = autoTLSConfig
		} else if server.SSLKey != "" && server.SSLCert != "" {
			tlsCFG = TLSConfig
		}

		if serverPortMap[server.Listen] == nil {
			serverPortMap[server.Listen] = &ServerPort{Port: server.Listen, TLSConfig: tlsCFG}
		} else if server.SSLAuto {
			serverPortMap[server.Listen].TLSConfig = tlsCFG
		}
	}
	return serverPortMap, nil
}

func GetTLSConfig(isAutoCert bool) (*tls.Config, error) {
	var cfg *tls.Config

	if isAutoCert {
		cfg = AutoCert.TLSConfig()
	} else {
		cfg = &tls.Config{}
		configs := GlobalConfigs
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
	}

	return cfg, nil
}

func GetClientRemoteAddr(r *http.Request) (string, string) {
	ip, port, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr))
	if err != nil {
		AddErrorLog(err)
		return "", ""
	} else {
		return ip, port
	}
}

func SetXForwardedRemotePort(r *http.Request) {
	_, port := GetClientRemoteAddr(r)
	xForwardedRemotePort := r.Header.Get(PortHeader)
	if len(xForwardedRemotePort) > 0 {
		xForwardedRemotePort += ", " + port
	} else {
		xForwardedRemotePort = port
	}
	r.Header.Set(PortHeader, xForwardedRemotePort)

}

func GetClientRealAddr(r *http.Request) (string, string) {
	configs := GlobalConfigs

	var ip string
	var port string

	if configs.Proxy && configs.MaxIpCount > 0 {
		ipsString := r.Header.Get(IpHeader)
		ipsStringArray := strings.Split(ipsString, ",")
		ipsCount := int16(len(ipsStringArray))
		if ipsCount >= configs.MaxIpCount {
			ip = ipsStringArray[ipsCount-configs.MaxIpCount]
		} else {
			ip = ""
		}

		portsString := r.Header.Get(PortHeader)
		portsStringArray := strings.Split(portsString, ", ")
		portsCount := int16(len(portsStringArray))
		if portsCount >= configs.MaxIpCount {
			port = portsStringArray[portsCount-configs.MaxIpCount]
		} else {
			port = ""
		}
	} else {
		ip, port = GetClientRemoteAddr(r)
	}
	return ip, port
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

func GetTargetService(host string, port uint16, url string) (*IService, error) {
	if GlobalServices[port] == nil ||
		GlobalServices[port][host] == nil ||
		len(GlobalServices[port][host]) == 0 {
		return nil, nil
	}
	services := GlobalServices[port][host]
	var targetService *IService
	for i := len(services) - 1; i >= 0; i-- {
		location := services[i].Location
		regString := location.Reg
		matched, err := regexp.MatchString(regString, url)
		if err != nil {
			return nil, err
		}
		if matched {
			targetService = &services[i]
			break
		}
	}
	return targetService, nil
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

func GetReqLimitByString(reqLimit []string) ([]*IReqLimit, error) {
	var reqLimitArr []*IReqLimit
	for _, item := range reqLimit {
		parameterError := errors.New("req_limit parameter error. req_limit=" + item)
		args := strings.Split(item, " ")
		argsLength := len(args)
		if argsLength < 2 || argsLength > 3 {
			return nil, parameterError
		}
		cacheNumberInt, err := strconv.Atoi(args[1])
		if err != nil {
			return nil, parameterError
		}
		CacheNumberUint64 := uint64(cacheNumberInt)

		reqLimitType := ReqLimitTypeStatic
		if len(args) > 2 && strings.TrimSpace(args[2]) == "ip" {
			reqLimitType = ReqLimitTypeIp
		}
		reqLimitTypeArr := strings.Split(strings.TrimSpace(args[0]), "/")
		if len(reqLimitTypeArr) != 2 {
			return nil, parameterError
		}

		countPerTimeInt, err := strconv.Atoi(reqLimitTypeArr[0])
		if err != nil {
			return nil, err
		}
		countPerTimeUint64 := uint64(countPerTimeInt)
		var timeNumber uint64 = 0

		duration, err := time.ParseDuration(reqLimitTypeArr[1])

		if err != nil {
			return nil, parameterError
		}

		timeNumber = uint64(duration.Milliseconds())

		reqLimit := &IReqLimit{
			Type:         reqLimitType,
			Time:         timeNumber,
			CountPerTime: countPerTimeUint64,
			CacheNumber:  CacheNumberUint64,
			Caches:       make(ICaches),
		}
		reqLimitArr = append(reqLimitArr, reqLimit)
	}
	return reqLimitArr, nil
}

func GetAbsPath(path string) (string, error) {
	if !filepath.IsAbs(path) {
		root, err := os.Getwd()
		if err != nil {
			return "", err
		}
		path = filepath.Join(root, path)
	}
	return path, nil
}

func IsDirValid(dirPath string) bool {
	fileInfo, err := os.Stat(dirPath)

	if err != nil {
		return false
	}

	if !fileInfo.IsDir() {
		return false
	}

	return true
}
