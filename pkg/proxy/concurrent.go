package proxy

import "time"

// 阻塞器
type blocker struct {
	c     chan uint8
	limit int64
}

// 当前时刻加入了一个处理中的请求
// 如果channel容量被占满了，那么send会被阻塞，后续的请求都会被停滞在这里
func (b *blocker) Add() {
	if b.limit <= 0 {
		return
	}
	b.c <- 1
}

// 当前通过了一个处理好的请求
// 调用后会从channel中receive一个数据，如果此时有被阻塞的请求，会在这之后继续处理
func (b *blocker) Pass() {
	if b.limit <= 0 {
		return
	}
	time.Sleep(time.Second)
	<-b.c
}

// 创建一个并发控制器
func NewConcurrentLimit(size int64) *blocker {
	b := new(blocker)
	b.c = make(chan uint8, size)
	b.limit = size
	return b
}
