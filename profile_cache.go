package sunshinemotion

import (
	"errors"
	"time"
)

type userSportResultCache struct {
	CacheTime      time.Time
	ExpireDuration time.Duration
	FetchFunction  func() (updated UserSportResult, err error)
	content        UserSportResult
}

// Get cached content
//
// if the cache is valid or has been updated successfully,
// error will be nil.

// if the cache has expired and an error occurred during the update process.
// it return the cached expired content and the updateError
func (cache *userSportResultCache) Get() (content UserSportResult, updateError error) {
	if cache.Expired() {
		err := cache.Update()
		if err != nil {
			return cache.content, errors.New("Get failed due to: " + err.Error())
		}
	}
	return cache.content, nil
}

// if return true, Update() will be called if needed
func (cache *userSportResultCache) Expired() bool {
	return time.Now().Before(cache.CacheTime.Add(cache.ExpireDuration))
}

// update cached content immediately
func (cache *userSportResultCache) Update() error {
	newContent, err := cache.FetchFunction()
	if err != nil {
		return errors.New("userInfoCache update failed: " + err.Error())
	}
	cache.CacheTime = time.Now()
	cache.content = newContent
	return nil
}
