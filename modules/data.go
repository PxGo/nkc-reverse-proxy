package modules

import (
	"errors"
	"time"
)

const (
	ReqLimitTypeStatic IReqLimitType = "static"
	ReqLimitTypeIp     IReqLimitType = "ip"
)

var GlobalServices IPortHostServices
var GlobalReqLimit []IReqLimit

type IService struct {
	Server   IServer
	Location ILocation
	Global   IGlobal
}
type IHostService map[string][]IService
type IPortHostServices map[uint16]IHostService

type IReqLimitType string
type IReqLimitChanData struct {
	LockChan IReqLimitLockChan
	Key      string
}

type ICacheChanData struct {
	Cache    *ICache
	ReqLimit IReqLimit
	LockChan IReqLimitLockChan
}

type IReqLimitChan chan IReqLimitChanData
type ICacheChan chan IReqLimitLockChan

type IReqLimitLockChan chan bool

type ICache struct {
	Time         time.Time
	Count        uint64
	WaitingCount uint64
	Chan         ICacheCheckerCoreChan
}

type ICacheCheckerCoreChanData struct {
	Cache    *ICache
	LockChan IReqLimitLockChan
}
type ICacheCheckerCoreChan chan ICacheCheckerCoreChanData

type ICaches map[string]*ICache

type IReqLimit struct {
	Type         IReqLimitType
	Time         uint64
	CountPerTime uint64
	CacheNumber  uint64
	Chan         IReqLimitChan
	Caches       ICaches
}

type ILocation struct {
	Reg          string
	Pass         []string
	Balance      string
	ReqLimit     *[]IReqLimit
	RedirectCode int
	RedirectUrl  string
}

type IServer struct {
	Listen   uint16
	Name     []string
	SSLKey   string
	SSLCert  string
	ReqLimit *[]IReqLimit
}

type IGlobal struct {
	ReqLimit *[]IReqLimit
}

// InitGlobalServices 在这里准备所有可能需要的数据
// 读取YAML文件内容
// 加载原始配置数据
// 转换部分字段数据
// 缓存进配置字段
func InitGlobalServices() error {
	if GlobalServices == nil {
		GlobalServices = make(IPortHostServices)
	}
	var err error
	GlobalReqLimit, err = GetReqLimitByString(GlobalConfigs.ReqLimit)

	iGlobal := IGlobal{
		ReqLimit: &GlobalReqLimit,
	}

	if err != nil {
		return err
	}
	for _, server := range GlobalConfigs.Servers {

		var services []IService

		if GlobalServices[server.Listen] == nil {
			GlobalServices[server.Listen] = make(IHostService)
		}
		serverReqLimit, err := GetReqLimitByString(server.ReqLimit)

		iServer := IServer{
			ReqLimit: &serverReqLimit,
		}

		if err != nil {
			return err
		}
		for _, location := range server.Location {
			locationReqLimit, err := GetReqLimitByString(location.ReqLimit)
			if err != nil {
				return err
			}
			iLocation := ILocation{
				Reg:          location.Reg,
				Pass:         location.Pass,
				Balance:      location.Balance,
				ReqLimit:     &locationReqLimit,
				RedirectCode: location.RedirectCode,
				RedirectUrl:  location.RedirectUrl,
			}
			services = append(services, IService{
				Global:   iGlobal,
				Server:   iServer,
				Location: iLocation,
			})
		}

		for _, name := range server.Name {
			if GlobalServices[server.Listen][name] == nil {
				GlobalServices[server.Listen][name] = services
			} else {
				return errors.New("duplicate domain name or port")
			}
		}
	}
	return nil
}
