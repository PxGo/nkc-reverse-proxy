package modules

import (
	"time"
)

func CreateReqLimitChan() IReqLimitChan {
	iChan := make(IReqLimitChan)
	return iChan
}

func RunReqLimitBroker(reqLimit *IReqLimit) {
	go func(args *IReqLimit) {
		CreateReqLimitChecker(reqLimit)
	}(reqLimit)
}

func RunCacheBroker(cache *ICache, lockChan IReqLimitLockChan) {
	go func() {
		cache.Chan <- ICacheCheckerCoreChanData{
			Cache:    cache,
			LockChan: lockChan,
		}
	}()
}

func InitCacheCheckerCore(reqLimit *IReqLimit) ICacheCheckerCoreChan {
	cacheCheckerCoreChan := make(ICacheCheckerCoreChan)
	go func(args *IReqLimit) {
		for {
			innerData := <-cacheCheckerCoreChan

			cache := innerData.Cache
			lockChan := innerData.LockChan

			now := time.Now()
			duration := now.Sub(cache.Time).Milliseconds()
			// 如果当前时段的请求数达到最大值且当前时段还有剩余时间
			if cache.Count >= reqLimit.CountPerTime && uint64(duration) < reqLimit.Time {
				// 条数超过了，需要定时再处理
				// 当前时段多余的时间，需要休眠
				timeLeft := reqLimit.Time - uint64(duration)
				time.Sleep(time.Duration(timeLeft) * time.Millisecond)
			}

			now = time.Now()
			duration = now.Sub(cache.Time).Milliseconds()

			// 如果已经到了下一个时间段则更新缓存时间和处理的请求数
			if uint64(duration) >= reqLimit.Time {
				// 已超过上一时段的时间，缓存重置
				cache.Time = now
				cache.Count = 0
			}

			cache.WaitingCount -= 1
			cache.Count += 1

			lockChan <- false
		}
	}(reqLimit)
	return cacheCheckerCoreChan
}

func CreateReqLimitChecker(reqLimit *IReqLimit) {

	for {
		reqLimitData := <-reqLimit.Chan

		lockChan := reqLimitData.LockChan
		key := reqLimitData.Key

		if reqLimit.Caches[key] == nil {
			iChan := InitCacheCheckerCore(reqLimit)
			reqLimit.Caches[key] = &ICache{
				Time:         time.Now(),
				Count:        0,
				WaitingCount: 0,
				Chan:         iChan,
			}
		}

		cache := reqLimit.Caches[key]

		if cache.WaitingCount >= reqLimit.CacheNumber {
			lockChan <- true
			continue
		}
		cache.WaitingCount += 1

		RunCacheBroker(cache, lockChan)
	}
}

func ReqLimitCheckerCore(reqLimit IReqLimit, key string) bool {
	lockChan := make(IReqLimitLockChan)

	reqLimit.Chan <- IReqLimitChanData{
		LockChan: lockChan,
		Key:      key,
	}

	limited := <-lockChan

	return limited
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
