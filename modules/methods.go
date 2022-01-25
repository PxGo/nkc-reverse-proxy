package modules

import (
	"crypto/tls"
	"errors"
	"github.com/dgryski/go-farm"
	"github.com/go-yaml/yaml"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
)

var proxyPassMap map[uint16]map[string]*ProxyPass
var httpReverseProxy *httputil.ReverseProxy
var httpsReverseProxy *httputil.ReverseProxy

func GetConfigsPath() (string, string, error) {
	root, err := os.Getwd()
	if err != nil {
		return "", "", err
	}
	filePath := path.Join(root, "configs.yaml")
	templateFilePath := path.Join(root, "configs.template.yaml")
	return filePath, templateFilePath, nil
}

func GetErrorLogPath() (string, error) {
	root, err := os.Getwd()
	if err != nil {
		return "", err
	}
	errorLogPath := path.Join(root, "error.log")
	return errorLogPath, nil
}

func GetConfigs() (*Configs, error) {
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
	var configs Configs
	err = yaml.Unmarshal(file, &configs)
	if err != nil {
		return nil, err
	}
	return &configs, nil
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

func GetProxyPassMap() (map[uint16]map[string]*ProxyPass, error) {
	if proxyPassMap != nil {
		return proxyPassMap, nil
	}
	proxyPass := make(map[uint16]map[string]*ProxyPass)
	configs, err := GetConfigs()
	if err != nil {
		return nil, err
	}
	for _, server := range configs.Servers {
		if proxyPass[server.Listen] == nil {
			proxyPass[server.Listen] = make(map[string]*ProxyPass)
		}
		for _, name := range server.Name {
			if proxyPass[server.Listen][name] == nil {
				var SocketIoPass []string
				var SocketIoBalance string
				if len(server.SocketIoPass) > 0 {
					SocketIoPass = server.SocketIoPass
				} else {
					SocketIoPass = server.Pass
				}
				if server.SocketIoBalance == "" {
					SocketIoBalance = server.balance
				} else {
					SocketIoBalance = server.SocketIoBalance
				}
				if server.RedirectUrl == "" && len(server.Pass) == 0 {
					return nil, errors.New("目标服务链接不能为空")
				}
				proxyPass[server.Listen][name] = &ProxyPass{
					Pass:            server.Pass,
					SocketIoPass:    SocketIoPass,
					balance:         server.balance,
					SocketIoBalance: SocketIoBalance,
					Redirect: RedirectInfo{
						Code: server.RedirectCode,
						Url:  server.RedirectUrl,
					},
				}
			}
		}
	}
	proxyPassMap = proxyPass
	return proxyPassMap, nil
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
	} else if passType == "random" {
		index = uint64(rand.Intn(passCount))
	} else {
		bytes := []byte(ip)
		hash := farm.Hash64(bytes)
		index = hash % uint64(passCount)
	}
	return pass[index]
}

func GetTargetPassInfo(req *http.Request, isHttps bool) (*url.URL, *RedirectInfo, error) {
	proxyPassMap, err := GetProxyPassMap()
	if err != nil {
		return nil, nil, err
	}
	hostInfo := strings.Split(req.Host, ":")
	if len(hostInfo) == 0 {
		return nil, nil, errors.New("客户端未指定 host")
	}
	host := hostInfo[0]
	polling := req.Header.Get("X-socket-io")
	isWS := polling == "polling"
	var port uint16
	var errInfo error
	if len(hostInfo) == 1 {
		if isHttps {
			port = 443
		} else {
			port = 80
		}
	} else {
		portInt, err := strconv.Atoi(hostInfo[1])
		if err != nil {
			return nil, nil, err
		}
		port = uint16(portInt)
		if errInfo != nil {
			return nil, nil, errInfo
		}
	}

	if err != nil {
		return nil, nil, err
	}
	var pass []string
	var passType string

	proxyPass := proxyPassMap[port][host]
	if proxyPass == nil {
		return nil, nil, nil
	}

	if isWS {
		pass = proxyPass.SocketIoPass
		passType = proxyPass.SocketIoBalance
	} else {
		pass = proxyPass.Pass
		passType = proxyPass.balance
	}

	var urlInfo *url.URL

	if len(pass) > 0 {
		ip := GetClientIP(req)
		targetUrlString := GetUrlByPassType(pass, passType, ip)
		var err error
		urlInfo, err = url.Parse(targetUrlString)
		if err != nil {
			return nil, nil, err
		}
	}

	return urlInfo, &proxyPass.Redirect, nil
}
