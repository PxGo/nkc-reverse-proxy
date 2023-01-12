package modules

import (
	"time"
)

type StoreData struct {
	Time  time.Time
	Count uint64
}

func CreateCacheChan() ICacheChan {
	iChan := make(ICacheChan)
	return iChan
}

func RunBroker(cache *ICache) {
	go func(args *ICache) {
		CreateReqLimitChecker(cache)
	}(cache)
}

func CreateReqLimitChecker(cache *ICache) {

	for {
		lockChan := <-cache.Chan

		now := time.Now()
		duration := now.Sub(cache.Time).Milliseconds()

		// 如果当前时段的请求数达到最大值且当前时段还有剩余时间
		if cache.Count >= cache.CountPerTime && uint64(duration) < cache.TimePerStage {
			// 条数超过了，需要定时再处理
			// 当前时段多余的时间，需要休眠
			timeLeft := cache.TimePerStage - uint64(duration)
			time.Sleep(time.Duration(timeLeft) * time.Millisecond)
		}

		now = time.Now()
		duration = now.Sub(cache.Time).Milliseconds()

		// 如果已经到了下一个时间段则更新缓存时间和处理的请求数
		if uint64(duration) >= cache.TimePerStage {
			// 已超过上一时段的时间，缓存重置
			cache.Time = now
			cache.Count = 0
		}
		cache.WaitingCount -= 1
		cache.Count += 1

		lockChan <- 1
	}
}

func ReqLimitCheckerCore(reqLimit IReqLimit, key string) bool {
	lockChan := make(ICacheLockChan)

	if reqLimit.Caches[key] == nil {
		iChan := CreateCacheChan()
		reqLimit.Caches[key] = &ICache{
			Time:         time.Now(),
			Count:        0,
			CountPerTime: reqLimit.CountPerTime,
			TimePerStage: reqLimit.Time,
			Chan:         iChan,
			WaitingCount: 0,
		}
		RunBroker(reqLimit.Caches[key])
	}

	cache := reqLimit.Caches[key]

	if cache.WaitingCount >= reqLimit.CacheNumber {
		// 当前时刻处理的请求数已经超过限制
		return true
	}

	cache.WaitingCount += 1

	cache.Chan <- lockChan

	<-lockChan

	return false
}

func ReqLimitChecker(reqLimit *[]IReqLimit, ip string) bool {
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
