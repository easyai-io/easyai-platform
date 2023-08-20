package localcache

import (
	"time"

	"github.com/patrickmn/go-cache"
)

var _cache *cache.Cache
var _kvCh chan *payload

type payload struct {
	key string
	val interface{}
	dur time.Duration
}

// nolint: gochecknoinits
func init() {
	_cache = cache.New(4*time.Hour, 30*time.Minute)
	_kvCh = make(chan *payload, 128)
	go func() {
		for kv := range _kvCh {
			_ = _cache.Add(kv.key, kv.val, kv.dur)
		}
	}()
}

// GetCache ...
func GetCache() *cache.Cache {
	return _cache
}

// SetKV async
func SetKV(prefix, key string, val interface{}, dur time.Duration) {
	_kvCh <- &payload{key: prefix + key, val: val, dur: dur}
}

// GetValue ...
func GetValue(prefix, key string) (interface{}, bool) {
	return _cache.Get(prefix + key)
}

// GetStringValue ...
func GetStringValue(prefix, key string) string {
	v, ok := _cache.Get(prefix + key)
	if ok {
		return v.(string)
	}
	return ""
}

// ----------

// LocalCache for local
type LocalCache struct {
	prefix string
}

// New LocalCache
func New(prefix string) *LocalCache {
	return &LocalCache{prefix: prefix}
}

// Exist key
func (l *LocalCache) Exist(key string) bool {
	_, ok := _cache.Get(l.prefix + key)
	return ok
}

// SetKVAsync ...
func (l *LocalCache) SetKVAsync(key string, val interface{}, dur time.Duration) {
	_kvCh <- &payload{key: l.prefix + key, val: val, dur: dur}
}

// SetKV ...
func (l *LocalCache) SetKV(key string, val interface{}, dur time.Duration) {
	_ = _cache.Add(l.prefix+key, val, dur)
}

// GetValue ...
func (l *LocalCache) GetValue(key string) (val interface{}, ok bool) {
	return _cache.Get(l.prefix + key)
}

// GetStringValue ...
func (l *LocalCache) GetStringValue(key string) (val string) {
	return GetStringValue(l.prefix, key)
}
