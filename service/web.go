package service

import (
	"fmt"
	"net/http"
	"strings"
)

// Web 简单的http处理
type Web struct {
	data  map[string]string
	store *Store
}

// NewWeb ..
func NewWeb(s *Store) *Web {
	w := new(Web)
	w.store = s
	return w
}
func (we *Web) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("新请求:" + r.Method + ":" + r.URL.Path + r.URL.RawQuery)
	if p := r.URL.Query().Get("get"); p != "" {
		w.Write([]byte(we.get(p)))
		return
	}
	if p := r.URL.Query().Get("set"); p != "" {
		val := strings.Split(p, ":")
		if len(val) == 2 {
			w.Write([]byte(we.set(val[0], val[1])))
			return
		}
		w.Write([]byte("错误的参数"))
		return
	}
	if p := r.URL.Query().Get("delete"); p != "" {
		w.Write([]byte(we.delete(p)))
		return
	}
	if p := r.URL.Query().Get("join"); p != "" {
		w.Write([]byte(we.join(p)))
		return
	}
}

// 获取一个key
func (we *Web) get(key string) string {
	return we.store.Get(key)
}

// 设置一个key
func (we *Web) set(key string, val string) string {
	if we.store.Set(key, val) == true {
		return "ok"
	}
	return "fail"
}

//删除一个key
func (we *Web) delete(key string) string {
	if we.store.Delete(key) == true {
		return "ok"
	}
	return "fail"
}

//接入一个新的服务节点
func (we *Web) join(address string) string {
	return "ok"
}
