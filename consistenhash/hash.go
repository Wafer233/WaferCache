package consistenhash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

type Hash func(data []byte) uint32

// 正常来说，先实现哈希一致性算法，事实上就是crc32这个包就解决了
// 然后为了防止数据倾斜，引入replicas，并且计算replica的哈希值
// 从小到大放入keys，并且引入一个hashmap指向真实的key
type Map struct {
	hash     Hash
	replicas int
	keys     []int // Sorted
	hashMap  map[int]string
}

func New(replicas int, hashFunc Hash) *Map {
	mapping := &Map{
		replicas: replicas,
		hash:     hashFunc,
		hashMap:  make(map[int]string),
	}

	if mapping.hash == nil {
		mapping.hash = crc32.ChecksumIEEE
	}
	return mapping
}

func (m *Map) Add(keys ...string) {

	for _, key := range keys {
		//函数允许传入 0 或 多个真实节点的名称。
		for i := 0; i < m.replicas; i++ {
			hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
			m.keys = append(m.keys, hash)
			m.hashMap[hash] = key
		}
	}
	sort.Ints(m.keys)
}

func (m *Map) Get(key string) string {
	if len(m.keys) == 0 {
		return ""
	}

	hash := int(m.hash([]byte(key)))

	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
	})

	return m.hashMap[m.keys[idx%len(m.keys)]]
}
