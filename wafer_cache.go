package WaferCache

import (
	"fmt"
	"log"
	"sync"
)

//	是 当前的key不在缓存的时候，要获取新的源
//
// 比较容易混淆的是 onEvicted，这个是lru淘汰之后如何处理旧的key，比如test里头的append进一个slice
type Loader interface {
	Load(key string) ([]byte, error)
}

type LoaderFunc func(key string) ([]byte, error)

func (function LoaderFunc) Load(key string) ([]byte, error) {
	return function(key)
}

type CacheNamespace struct {
	name   string
	loader Loader
	cache  safeCache
}

// 这里涉及到新建namesapce和读取namespace，所以也要加锁；
// 然后由于读多写少（因为一般不会频繁创建新的nameSpace），所以用rwmutex
// nsMu = Namespace Mutex
var nsMu sync.RWMutex
var nameSpaces = make(map[string]*CacheNamespace)

func NewNameSpace(name string, maxBytes int64, loader Loader) *CacheNamespace {
	if loader == nil {
		panic("loader没有设置")
	}

	nsMu.Lock()
	defer nsMu.Unlock()

	newSpace := &CacheNamespace{
		name:   name,
		loader: loader,
		cache:  safeCache{maxBytes: maxBytes},
	}

	nameSpaces[name] = newSpace
	return newSpace
}

func GetNameSpace(name string) *CacheNamespace {
	nsMu.RLock()
	defer nsMu.RUnlock()

	curSpace := nameSpaces[name]

	return curSpace
}

// Get value for a key from cache
func (ns *CacheNamespace) Get(key string) (ValueView, error) {
	if key == "" {
		return ValueView{}, fmt.Errorf("需要key")
	}

	// 流程 1 ：从 cache 中查找缓存，如果存在则返回缓存值。
	if view, exist := ns.cache.get(key); exist {
		log.Println("命中缓存")
		return view, nil
	}

	// 流程3：缓存不存在，则调用 load 方法，load 调用 getLocally (单机并发)
	return ns.load(key)
}

func (ns *CacheNamespace) load(key string) (ValueView, error) {
	return ns.loadLocally(key)
}

func (ns *CacheNamespace) loadLocally(key string) (ValueView, error) {
	bytes, err := ns.loader.Load(key)
	if err != nil {
		return ValueView{}, err

	}
	value := ValueView{b: cloneBytes(bytes)}

	ns.setCache(key, value)
	return value, nil
}

func (ns *CacheNamespace) setCache(key string, value ValueView) {
	ns.cache.set(key, value)
}
