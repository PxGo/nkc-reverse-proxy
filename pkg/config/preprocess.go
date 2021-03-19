package config

import (
	"crypto/tls"
	"fmt"
	"log"
)

func preprocess(conf *Profile) {
	set_NoResponsePage(conf)
	set_SSL(conf)
}

// 检查局部NoResponsePage配置缺省，填充全局NoResponsePage配置
func set_NoResponsePage(conf *Profile) {
	for i, serve := range conf.Servers {
		if serve.NoResponsePage == "" {
			conf.Servers[i].NoResponsePage = conf.NoResponsePage
		}
	}
}

// 检查和设置每个Server配置的SSL配置项
// 检查局部SSL配置缺省，填充全局SSL配置
func set_SSL(conf *Profile) {
	defaultSSL := conf.SSL
	isMissDefaultSSL := defaultSSL.Cert == "" || defaultSSL.Key == ""
	for i, serve := range conf.Servers {
		if serve.Https {
			if serve.SSL.Cert == "" || serve.SSL.Key == "" {
				if isMissDefaultSSL {
					log.Fatal("有代理规则配置了https访问，但未配置局部或全局SSL配置项")
				} else {
					conf.Servers[i].SSL.Cert = defaultSSL.Cert
					conf.Servers[i].SSL.Key = defaultSSL.Key
				}
			} else {
				cacheSSLCertificate(serve.SSL.Cert, serve.SSL.Key)
			}
		}
	}
	if !isMissDefaultSSL {
		cacheSSLCertificate(defaultSSL.Cert, defaultSSL.Key)
	}
}

var CertificateCaches = map[string]*tls.Certificate{}

// 预加载SSL证书
func cacheSSLCertificate(certFile string, keyFile string) {
	if CertificateCaches[certFile] != nil {
		return
	}
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		fmt.Println("读取SSL证书文件失败")
		panic(err)
	}
	CertificateCaches[certFile] = &cert
}
