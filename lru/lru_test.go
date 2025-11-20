package lru

import (
	"reflect"
	"testing"
)

// 设置一个testVal 是一个Value类型

type String string

func (testVal String) Len() int {
	return len(testVal)
}

func TestSet(t *testing.T) {
	//无限缓存
	lru := New(int64(0), nil)
	lru.Set("key", String("1"))
	lru.Set("key", String("111"))

	if lru.curBytes != int64(len("key")+len("111")) {
		t.Fatal("想要6，但是获得了:", lru.curBytes)
	}
}

func TestGet(t *testing.T) {
	lru := New(int64(0), nil)
	lru.Set("key1", String("1234"))

	if v, exist := lru.Get("key1"); !exist || string(v.(String)) != "1234" {
		t.Fatalf("查询key1失败")
	}
	if _, exist := lru.Get("key2"); exist {
		t.Fatalf("查询到不该有的key2")
	}
}

func TestRemoveOldest(t *testing.T) {
	k1, k2, k3 := "key1", "key2", "key3"
	v1, v2, v3 := "value1", "value2", "v3"

	curByte := len(k1 + k2 + v1 + v2)
	lru := New(int64(curByte), nil)
	lru.Set(k1, String(v1))
	lru.Set(k2, String(v2))
	lru.Set(k3, String(v3))

	if _, ok := lru.Get("key1"); ok || lru.Len() != 2 {
		t.Fatalf("移除key1失败")
	}

}

func TestOnEvicted(t *testing.T) {
	keys := make([]string, 0)
	callback := func(key string, value Value) {
		keys = append(keys, key)
	}
	lru := New(int64(10), callback)
	lru.Set("key1", String("123456"))
	lru.Set("k2", String("k2"))
	lru.Set("k3", String("k3"))
	lru.Set("k4", String("k4"))

	expect := []string{"key1", "k2"}

	if !reflect.DeepEqual(expect, keys) {
		t.Fatalf("Call OnEvicted failed, expect keys equals to %s", expect)
	}
}
