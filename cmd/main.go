package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Wafer233/WaferCache/WaferCache"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

// http://localhost:9999/wafer_cache/scores/Tom
func main() {
	WaferCache.NewNameSpace("scores", 2<<10, WaferCache.LoaderFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] 搜索key", key)
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s 不存在", key)
		}))

	addr := "localhost:9999"
	peers := WaferCache.NewHTTPPool(addr)
	log.Println("WaferCache 运行在以下地址", addr)
	log.Fatal(http.ListenAndServe(addr, peers))
}
