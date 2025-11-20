package WaferCache

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

const defaultBasePath = "/wafer_cache/"

// 约定访问路径格式 /<basepath>/<namespace>/<key>
type HTTPPool struct {
	self     string
	basePath string
}

func NewHTTPPool(self string) *HTTPPool {
	return &HTTPPool{
		self:     self,
		basePath: defaultBasePath,
	}
}

func (pool *HTTPPool) Log(format string, v ...interface{}) {
	log.Printf("[Server %s] %s", pool.self, fmt.Sprintf(format, v...))
}

func (pool *HTTPPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 判断不合法
	if !strings.HasPrefix(r.URL.Path, pool.basePath) {
		panic("basepath不合法: " + r.URL.Path)
	}
	pool.Log("%s %s", r.Method, r.URL.Path)

	// /<basepath>/<groupname>/<key> required
	parts := strings.SplitN(r.URL.Path[len(pool.basePath):], "/", 2)
	if len(parts) != 2 {
		http.Error(w, "请求错误，请按照/<basepath>/<namespace>/<key> 请求", http.StatusBadRequest)
		return
	}

	curName := parts[0]
	key := parts[1]

	nameSpace := GetNameSpace(curName)
	if nameSpace == nil {
		http.Error(w, "没有当前的namespace: "+curName, http.StatusNotFound)
		return
	}

	view, err := nameSpace.Get(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//流式的字节数据
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(view.ByteSlice())
}
