package WaferCache

import (
	"fmt"
	"log"
	"reflect"
	"testing"
)

func TestLoader(t *testing.T) {
	var function Loader = LoaderFunc(func(key string) ([]byte, error) {
		return []byte("自定义loader的加载的key:" + key), nil
	})

	expect := []byte("自定义loader的加载的key:wafer_key")
	if v, _ := function.Load("wafer_key"); !reflect.DeepEqual(v, expect) {
		t.Fatal("callback failed")
	}
}

func TestGetNameSpace(t *testing.T) {
	name := "scores"
	maxBytes := int64(2 << 10)
	loader := LoaderFunc(func(key string) ([]byte, error) {
		return []byte{}, nil
	})

	NewNameSpace(name, maxBytes, loader)

	if space := GetNameSpace(name); space == nil || space.name != name {
		t.Fatalf("当前namespace不存在")
	}

	if space := GetNameSpace(name + "111"); space != nil {
		t.Fatalf("查询到了不该有的namespace")
	}
}

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func TestGet(t *testing.T) {
	loadCounts := make(map[string]int, len(db))

	name := "scores"
	maxBytes := int64(2 << 10)
	loader := LoaderFunc(func(key string) ([]byte, error) {

		log.Println("[SlowDB] search key", key)
		if value, exist := db[key]; exist {

			if _, exist := loadCounts[key]; !exist {
				loadCounts[key] = 0
			}
			loadCounts[key]++
			return []byte(value), nil
		}
		return nil, fmt.Errorf("%s not exist", key)
	})

	curNameSpace := NewNameSpace(name, maxBytes, loader)

	for key, value := range db {

		// 正常来说要能够读到key
		if view, err := curNameSpace.Get(key); err != nil || view.String() != value {
			t.Fatal("failed to get value of Tom")
		}

		// 而且不能加载两次，因为第二次就可以直接从cache中获取，而不死需要从db中获取
		if _, err := curNameSpace.Get(key); err != nil || loadCounts[key] > 1 {
			t.Fatalf("cache %s miss", key)
		}
	}

	//不能够获取没有缓存的key
	if view, err := curNameSpace.Get("unknown"); err == nil {
		t.Fatalf("the value of unknow should be empty, but %s got", view)
	}
}
