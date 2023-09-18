package datastruct

import (
	"math"
	"sync"
	"sync/atomic"
)

type ConcurrentMap struct {
	table []*shard
	count int32
}

// 构造器方法
func NewConcurrentMap(cap int) *ConcurrentMap {
	// 计算得到分片的个数
	shardCount := getShardCount(cap)
	table := make([]*shard, 0, shardCount)
	for i := 0; i < shardCount; i++ {
		table = append(table, &shard{
			data: make(map[string]interface{}),
		})
	}
	return &ConcurrentMap{
		table: table,
	}
}

func (c *ConcurrentMap) Get(key string) (interface{}, bool) {
	shard := c.getShard(key)
	shard.lock.RLock()
	defer shard.lock.RUnlock()
	v, ok := shard.data[key]
	return v, ok
}

func (c *ConcurrentMap) Put(key string, val interface{}) int {
	shard := c.getShard(key)
	shard.lock.Lock()
	defer shard.lock.Unlock()
	var reply int
	if _, ok := shard.data[key]; !ok {
		reply = 1
		c.addCount()
	}
	shard.data[key] = val
	return reply
}

func (c *ConcurrentMap) getShardIndex(key string) uint32 {
	hash := fvn32(key)
	return hash & uint32(len(c.table)-1)
}

func (c *ConcurrentMap) getShard(key string) *shard {
	if c.table == nil {
		return nil
	}
	shadIndex := c.getShardIndex(key)
	return c.table[shadIndex]
}

func (c *ConcurrentMap) addCount() {
	atomic.AddInt32(&c.count, 1)
}

type shard struct {
	lock sync.RWMutex
	data map[string]interface{}
}

func getShardCount(cap int) int {
	if cap <= 16 {
		return 16
	}

	n := cap - 1
	n |= n >> 1
	n |= n >> 2
	n |= n >> 4
	n |= n >> 8
	n |= n >> 16
	if n < 0 {
		return int(math.Pow(2, 30))
	}

	return n + 1
}
