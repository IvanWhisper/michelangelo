package modulo

import "michelangelo/log"

// 取模
type Sharding struct {
	Len uint32
}

// 获取分片Id
func (s *Sharding) GetShardIndex(id uint64) uint64 {
	return id % uint64(s.Len)
}

// 新建实例
func New(len uint32) *Sharding {
	m := &Sharding{Len: len}
	log.Info("Create Len%s ")
	return m
}
