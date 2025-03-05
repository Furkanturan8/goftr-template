package cache

import "errors"

var (
	ErrCacheNotInitialized = errors.New("cache is not initialized")
	ErrKeyNotFound         = errors.New("key not found in cache")
	ErrInvalidValue        = errors.New("invalid value type")
)
