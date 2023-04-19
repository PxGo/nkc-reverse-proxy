package modules

import (
	"sync"
	"time"
)

const (
	ReqLimitTypeStatic IReqLimitType = "static"
	ReqLimitTypeIp     IReqLimitType = "ip"
)

type IReqLimitType string

type ICaches map[string]*ICache

type IReqLimit struct {
	Mu           sync.Mutex
	Type         IReqLimitType
	Time         uint64
	CountPerTime uint64
	CacheNumber  uint64
	Caches       ICaches
}

func (own *IReqLimit) ClearCacheByKey(key string) {
	own.LockCache()
	delete(own.Caches, key)
	own.UnlockCache()
}

func (own *IReqLimit) LockCache() {
	own.Mu.Lock()
}

func (own *IReqLimit) UnlockCache() {
	own.Mu.Unlock()
}

type IReqLimitLockChan chan bool

type ICache struct {
	Lock         *sync.Mutex
	LoopLock     *sync.Mutex
	Time         time.Time
	Count        uint64
	WaitingCount uint64

	reqLimitCounterPerTime uint64
	reqLimitTime           uint64
	reqLimitCacheNumber    uint64
}

func (cache *ICache) LockWaitingCountToMinusOne() {
	cache.LockWaitingCount()
	if cache.WaitingCount > 0 {
		cache.WaitingCount -= 1
	}
	cache.UnlockWaitingCount()
}

func (cache *ICache) LockWaitingCount() {
	cache.Lock.Lock()
}

func (cache *ICache) UnlockWaitingCount() {
	cache.Lock.Unlock()
}

func (cache *ICache) LockLoop() {
	cache.LoopLock.Lock()
}

func (cache *ICache) UnlockLoop() {
	cache.LoopLock.Unlock()
}

func (cache *ICache) IsCacheFull() bool {
	return cache.WaitingCount >= cache.reqLimitCacheNumber
}

func (cache *ICache) CreateLockChan() IReqLimitLockChan {

	lockChan := make(IReqLimitLockChan)

	go func() {

		cache.LoopLock.Lock()

		for {
			now := time.Now()

			duration := now.Sub(cache.Time).Milliseconds()
			// 如果当前时段的请求数达到最大值且当前时段还有剩余时间
			if cache.Count >= cache.reqLimitCounterPerTime && uint64(duration) < cache.reqLimitTime {
				// 条数超过了，需要定时再处理
				// 当前时段多余的时间，需要休眠
				timeLeft := cache.reqLimitTime - uint64(duration)
				time.Sleep(time.Duration(timeLeft) * time.Millisecond)

				continue
			}

			// 如果已经到了下一个时间段则更新缓存时间和处理的请求数
			if uint64(duration) >= cache.reqLimitTime {
				// 已超过上一时段的时间，缓存重置
				cache.Time = now
				cache.Count = 0
			}

			cache.Count += 1

			cache.LockWaitingCountToMinusOne()

			lockChan <- false

			break
		}

		cache.LoopLock.Unlock()
	}()

	return lockChan
}

func ReqLimitCheckerCore(reqLimit *IReqLimit, key string) bool {

	reqLimit.LockCache()

	if reqLimit.Caches[key] == nil {

		var loopLock sync.Mutex
		var lock sync.Mutex

		cache := &ICache{
			reqLimitCounterPerTime: reqLimit.CountPerTime,
			reqLimitTime:           reqLimit.Time,
			reqLimitCacheNumber:    reqLimit.CacheNumber,

			Time:         time.Now(),
			Count:        0,
			WaitingCount: 0,
			LoopLock:     &loopLock,
			Lock:         &lock,
		}

		reqLimit.Caches[key] = cache

		go func() {
			var timeout = int64(cache.reqLimitTime + 1000*60*60)
			for {
				time.Sleep(time.Millisecond * time.Duration(timeout))
				now := time.Now()
				duration := now.Sub(cache.Time)
				if duration.Milliseconds() > timeout {
					reqLimit.ClearCacheByKey(key)
					return
				} else {
					continue
				}
			}
		}()
	}

	cache := reqLimit.Caches[key]

	reqLimit.UnlockCache()

	cache.LockWaitingCount()

	if cache.IsCacheFull() {

		cache.UnlockWaitingCount()

		return true
	}

	cache.WaitingCount += 1

	cache.UnlockWaitingCount()

	lockChan := cache.CreateLockChan()

	limited := <-lockChan

	return limited

}

func ReqLimitChecker(reqLimit *[]*IReqLimit, ip string) bool {
	for _, item := range *reqLimit {
		key := "static"
		if item.Type == ReqLimitTypeIp {
			key = ip
		}
		limited := ReqLimitCheckerCore(item, key)
		if limited {
			return limited
		}
	}
	return false
}
