package sunshinemotion

import (
	"errors"
	"time"
)

type sportResultCache struct {
	// update cached content immediately.
	// 由缓存使用者提供。在Update函数中，调用cache.Put(content, time)更新缓存。
	Update         func() (err error)
	ExpireDuration time.Duration // 过期时间
	content        UserSportResult
	cacheTime      time.Time
}

// Get cached content (maybe expired)
func (cache *sportResultCache) Get() (content UserSportResult, cacheTime time.Time) {
	return cache.content, cache.cacheTime
}

// set cached content
func (cache *sportResultCache) Put(content UserSportResult, cacheTime time.Time) {
	cache.cacheTime = cacheTime
	cache.content = content
}

// if return true, Update() will be called if needed
func (cache *sportResultCache) Expired() bool {
	return time.Now().Before(cache.cacheTime.Add(cache.ExpireDuration))
}

// Get valid content
//
// if the cache is valid or has been updated successfully,
// error will be nil.
//
// if the cache has expired and an error occurred during the update process.
// it return the cached expired content and the updateError
func (cache *sportResultCache) GetValid() (content UserSportResult, cacheTime time.Time, updateError error) {
	if cache.Expired() {
		if cache.Update == nil {
			return cache.content, cache.cacheTime, errors.New("cache update function not provided")
		}
		err := cache.Update()
		if err != nil {
			return cache.content, cache.cacheTime, errors.New("Get failed due to: " + err.Error())
		}
	}
	return cache.content, cache.cacheTime, nil
}
