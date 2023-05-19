package modules

import (
	"golang.org/x/crypto/acme/autocert"
	"os"
	"path/filepath"
)

var AutoCert *autocert.Manager

func InitAutoCert() error {
	cacheDirPath, err := GetAutoCertCacheDirPath()
	if err != nil {
		return nil
	}
	hostWhitelist := GetAutoCertHostWhitelist()
	AutoCert = &autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(hostWhitelist...),
		Cache:      autocert.DirCache(cacheDirPath),
	}
	return nil
}

func GetAutoCertCacheDirPath() (string, error) {
	cacheDir := ".cache"

	root, err := os.Getwd()
	if err != nil {
		return "", err
	}

	cacheDirPath := filepath.Join(root, cacheDir)
	return cacheDirPath, nil
}

func GetAutoCertHostWhitelist() []string {
	var host []string
	configs := GlobalConfigs
	for _, server := range configs.Servers {
		host = append(host, server.Name...)
	}
	return host
}

func AutoCertIsEnabled() bool {
	hostWhitelist := GetAutoCertHostWhitelist()
	return len(hostWhitelist) > 0
}
