package modules

import (
	"sync"
	"time"
)

type IIpLimitCacheStatus string

type IIpLimitCache struct {
	Lock         *sync.Mutex
	Banned       bool
	StartTime    time.Time
	VisitedCount uint64

	UnitTime         uint64
	CountPerUnitTime uint64
	BanTime          uint64
}

type IIpLimitCaches map[string]*IIpLimitCache

type IIpLimit struct {
	Mu               sync.Mutex
	UnitTime         uint64
	CountPerUnitTime uint64
	BanTime          uint64
	Caches           IIpLimitCaches
}

type IIpLimitLockChan chan bool

func (own *IIpLimit) Lock() {
	own.Mu.Lock()
}

func (own *IIpLimit) Unlock() {
	own.Mu.Unlock()
}

func (own *IIpLimit) ClearCacheByKey(key string) {
	own.Lock()
	delete(own.Caches, key)
	own.Unlock()
}

func (cache *IIpLimitCache) LockCache() {
	cache.Lock.Lock()
}

func (cache *IIpLimitCache) UnLockCache() {
	cache.Lock.Unlock()
}

func IpLimitCheckerCore(ipLimit *IIpLimit, ip string) bool {
	ipLimit.Lock()

	cache, ok := ipLimit.Caches[ip]
	if !ok {
		cache = &IIpLimitCache{
			Lock:             &sync.Mutex{},
			Banned:           false,
			StartTime:        time.Now(),
			VisitedCount:     0,
			UnitTime:         ipLimit.UnitTime,
			CountPerUnitTime: ipLimit.CountPerUnitTime,
			BanTime:          ipLimit.BanTime,
		}
		ipLimit.Caches[ip] = cache

		go func() {
			timeout := int64(cache.UnitTime + 1000*60*60)
			for {
				time.Sleep(time.Millisecond * time.Duration(timeout))
				now := time.Now()
				duration := now.Sub(cache.StartTime)
				if duration.Milliseconds() > timeout {
					ipLimit.ClearCacheByKey(ip)
					return
				} else {
					continue
				}
			}
		}()
	}

	cache = ipLimit.Caches[ip]

	ipLimit.Unlock()

	cache.LockCache()

	now := time.Now()
	duration := now.Sub(cache.StartTime)
	durationUint64 := uint64(duration.Milliseconds())

	if cache.Banned {
		// 已被封禁，判断是否可以解禁
		if durationUint64 > cache.BanTime {
			// 超过了封禁时间
			cache.Banned = false
			cache.StartTime = time.Now()
			cache.VisitedCount = 0
		}
	} else {
		if durationUint64 > cache.UnitTime {
			// 超过了一个时间段
			cache.StartTime = time.Now()
			cache.VisitedCount = 0
		} else {
			if cache.VisitedCount >= cache.CountPerUnitTime {
				// 访问此数超过了限制
				cache.StartTime = time.Now()
				cache.Banned = true
			} else {
				cache.VisitedCount += 1
			}
		}
	}

	return cache.Banned
}

func IpLimitChecker(ipLimit *[]*IIpLimit, ip string) bool {
	for _, item := range *ipLimit {
		limited := IpLimitCheckerCore(item, ip)
		if limited {
			return limited
		}
	}
	return false
}
