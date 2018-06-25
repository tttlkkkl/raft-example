package service

import (
	"sync"

	"github.com/hashicorp/raft"
)

// Store 简易存储结构
type Store struct {
	data map[string]string
	mux  *sync.Mutex
	raft *raft.Raft
}

// NewStore 实例化存储库
func NewStore(raft *raft.Raft) *Store {
	s := new(Store)
	s.data = make(map[string]string)
	s.mux = new(sync.Mutex)
	s.raft = raft
	return s
}

// Get 从存储库中获取值
func (s *Store) Get(key string) string {
	if key == "" {
		return ""
	}
	s.mux.Lock()
	defer s.mux.Unlock()
	val, ok := s.data[key]
	if ok == true {
		return val
	}
	return ""
}

// Set 存储值
func (s *Store) Set(key string, val string) bool {
	if key == "" {
		return false
	}
	s.mux.Lock()
	defer s.mux.Unlock()
	s.data[key] = val
	return true
}

// Delete 删除值
func (s *Store) Delete(key string) bool {
	if key == "" {
		return false
	}
	s.mux.Lock()
	defer s.mux.Unlock()
	delete(s.data, key)
	return true
}
