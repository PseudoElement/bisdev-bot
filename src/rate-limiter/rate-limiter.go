package ratelimiter

import "time"

const DELAY = 10 * time.Second
const RPS_LIMIT = 10

type reqFromClient struct {
	firstMsgTimestamp int64
	msgCount          int
}

type RateLimiter struct {
	cache map[int64]reqFromClient
}

func NewRateLimiter() *RateLimiter {
	return &RateLimiter{cache: make(map[int64]reqFromClient)}
}

func (this *RateLimiter) CacheReq(userId int64) {
	_, ok := this.cache[userId]
	if ok {
		oldVal := this.cache[userId]
		oldVal.msgCount++
		this.cache[userId] = oldVal
	} else {
		this.cache[userId] = reqFromClient{
			firstMsgTimestamp: time.Now().UnixMilli(),
			msgCount:          1,
		}
	}

	go func() {
		time.Sleep(DELAY)
		delete(this.cache, userId)
	}()
}

func (this *RateLimiter) CheckOnSpam(userId int64) bool {
	req, ok := this.cache[userId]
	return ok && req.msgCount > RPS_LIMIT
}
