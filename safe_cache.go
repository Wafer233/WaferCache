package WaferCache

import (
	"sync"

	"github.com/Wafer233/WaferCache/WaferCache/lru"
)

// 这里的cache其实就是给lru上锁
type safeCache struct {
	mu       sync.Mutex
	lru      *lru.Cache
	maxBytes int64
}

// 上锁，并发
func (cache *safeCache) set(key string, value ValueView) {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	// 如果没有
	if cache.lru == nil {
		cache.lru = lru.New(cache.maxBytes, nil)
	}

	// 新增或者更新
	cache.lru.Set(key, value)
}

func (cache *safeCache) get(key string) (ValueView, bool) {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	// 空
	if cache.lru == nil {
		return ValueView{}, false
	}

	//存在
	if value, exist := cache.lru.Get(key); exist {
		return value.(ValueView), true
	}

	//不存在
	return ValueView{}, false
}
