package lru

import (
	"container/list"
)

//首先就要引入淘汰机制
//正常来说LRU 对比FIFO和LFU肯定是最好的，前者会导致缓存一直更新，后者会导致内存负载变高

type Cache struct {
	maxBytes         int64
	curBytes         int64
	doublyLinkedList *list.List
	cache            map[string]*list.Element //加速定位key在doubly linkedlist中的位置 O(1)
	OnEvicted        func(key string, value Value)
}

type pair struct {
	key   string
	value Value
}

type Value interface {
	Len() int
}

// New 初始化Cache
func New(maxBytes int64, onEvicted func(key string, value Value)) *Cache {
	return &Cache{
		maxBytes:         maxBytes,
		curBytes:         0,
		doublyLinkedList: list.New(),
		cache:            make(map[string]*list.Element),
		OnEvicted:        onEvicted,
	}
}

// Set 新增缓存（修改）
func (cache *Cache) Set(key string, value Value) {

	//先看key是不是存在 -> 不存在 我要在最前面插入一个pair
	if element, exists := cache.cache[key]; !exists {

		newKV := pair{
			key:   key,
			value: value,
		}
		element = cache.doublyLinkedList.PushFront(&newKV)
		cache.cache[key] = element
		addByte := int64(len(key)) + int64(value.Len())
		cache.curBytes += addByte
	} else {
		// 如果存在 ->挪到最前面，大小不变 (更新value) 变换大小
		cache.doublyLinkedList.MoveToFront(element)
		kv := element.Value.(*pair)
		oldSize := kv.value.Len()
		newSize := value.Len()

		kv.value = value
		cache.curBytes += int64(newSize - oldSize)
	}

	// 删除最老的
	// cache.maxBytes ==0 -> 无限容量
	for cache.maxBytes != 0 && cache.curBytes > cache.maxBytes {
		cache.RemoveOldest()
	}
}

// Get 查找缓存
func (cache *Cache) Get(key string) (Value, bool) {

	// 查到了
	if element, exists := cache.cache[key]; exists {
		cache.doublyLinkedList.MoveToFront(element)
		kv := element.Value.(*pair)
		return kv.value, true
	}

	// 没查到
	return nil, false
}

// Delete 删除缓存 （满了）
func (cache *Cache) RemoveOldest() {
	element := cache.doublyLinkedList.Back()

	//删除最后一个
	if element != nil {
		cache.doublyLinkedList.Remove(element)
		kv := element.Value.(*pair)
		delete(cache.cache, kv.key)
		delByte := int64(len(kv.key)) + int64(kv.value.Len())
		cache.curBytes -= delByte

		//删除后调用callback函数
		if cache.OnEvicted != nil {
			cache.OnEvicted(kv.key, kv.value)
		}
	}
}

// Len 为了方便测试
func (cache *Cache) Len() int {
	return cache.doublyLinkedList.Len()
}
