package config

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"runtime/debug"
)

func Parse(file string) Profile {
	config := Profile{}
	err := json.Unmarshal(loadFile(file), &config)
	if err != nil {
		fmt.Printf("%s: %s\n%s\n", "配置文件读取出错", err, debug.Stack())
		os.Exit(1)
	}
	preprocess(&config)
	return config
}

func loadFile(file string) []byte {
	content, err := os.ReadFile(file)
	if err != nil {
		fmt.Println("无法读取文件: " + file)
		log.Fatal(err)
	}
	return content
}

func GetCertificate(name string) (*tls.Certificate, error) {
	if CertificateCaches[name] == nil {
		return nil, errors.New("在缓存中没有找到此SSL证书对象")
	}
	return CertificateCaches[name], nil
}
