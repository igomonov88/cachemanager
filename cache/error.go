package cache

type cacheError string

func (err cacheError) Error() string {
	return string(err)
}

const (
	ErrOverCacheLimit     = cacheError("go over the cache limit")
	ErrNoValueForGivenKey = cacheError("there is no value for given key")
)
