package config

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
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
	content, err := os.ReadFile(GetAbsPath(file))
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

func GetAbsPath(path string) string {
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	if filepath.IsAbs(path) {
		return path
	} else {
		return filepath.Join(pwd, path)
	}
}

func PathExist(_path string) bool {
	_, err := os.Stat(_path)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}
