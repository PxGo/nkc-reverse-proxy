package modules

import (
	"sync"
	"time"
)

func ReqLimitCheckerCore(reqLimit *IReqLimit, key string) bool {

	reqLimit.Mu.Lock()
	if reqLimit.Caches[key] == nil {

		var loopLock sync.Mutex
		var lock sync.Mutex

		cond := sync.NewCond(&loopLock)

		cache := &ICache{
			Time:         time.Now(),
			Count:        0,
			WaitingCount: 0,

			LoopLock:    &loopLock,
			Lock:        &lock,
			Cond:        cond,
			LockChanArr: []IReqLimitLockChan{},
			MarkClear:   false,
		}

		reqLimit.Caches[key] = cache

		// 执行任务
		go func() {
			for {
				cache.LoopLock.Lock()

				// 如果队列里没有请求，则等待
				for len(cache.LockChanArr) == 0 {
					cache.Cond.Wait()
				}

				cache.LoopLock.Unlock()

				if cache.MarkClear {
					return
				}

				for {

					if len(cache.LockChanArr) == 0 {
						continue
					}

					now := time.Now()

					duration := now.Sub(cache.Time).Milliseconds()
					// 如果当前时段的请求数达到最大值且当前时段还有剩余时间
					if cache.Count >= reqLimit.CountPerTime && uint64(duration) < reqLimit.Time {
						// 条数超过了，需要定时再处理
						// 当前时段多余的时间，需要休眠
						timeLeft := reqLimit.Time - uint64(duration)
						time.Sleep(time.Duration(timeLeft) * time.Millisecond)
						continue
					}

					cache.Lock.Lock()

					// 如果已经到了下一个时间段则更新缓存时间和处理的请求数
					if uint64(duration) >= reqLimit.Time {
						// 已超过上一时段的时间，缓存重置
						cache.Time = now
						cache.Count = 0
					}

					cache.WaitingCount -= 1
					cache.Count += 1

					lockChan := cache.LockChanArr[0]
					cache.LockChanArr = cache.LockChanArr[1:]

					cache.Lock.Unlock()

					lockChan <- false

				}
			}
		}()

		var timeout = int64(reqLimit.Time * 10)
		// 资源回收
		go func() {
			for {
				time.Sleep(time.Millisecond * time.Duration(timeout))
				now := time.Now()
				duration := now.Sub(reqLimit.Caches[key].Time)
				if duration.Milliseconds() > timeout {
					//fmt.Println("开始清理cache:", key)
					cache.ClearGoroutine()
					reqLimit.ClearCacheByKey(key)
					return
				} else {
					continue
				}
			}
		}()
	}
	reqLimit.Mu.Unlock()

	cache := reqLimit.Caches[key]

	cache.Lock.Lock()

	if cache.WaitingCount >= reqLimit.CacheNumber {
		cache.Lock.Unlock()
		return true
	}

	cache.WaitingCount += 1

	lockChan := make(IReqLimitLockChan)

	cache.LockChanArr = append(cache.LockChanArr, lockChan)

	cache.Lock.Unlock()

	cache.Cond.Broadcast()

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
