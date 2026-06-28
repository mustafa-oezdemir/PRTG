package prtg

import (
	"sync"
	"time"
)

type cacheItem struct {
	data   []byte
	expiry time.Time
}

type Api struct {
	baseURL   string
	apiKey    string
	timeout   time.Duration
	cacheTime time.Duration
	cache     map[string]cacheItem
	cacheMu   sync.RWMutex
}
