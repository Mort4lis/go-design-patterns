package concurrency

import (
	"crypto/sha1"
	"sync"
)

type Shard struct {
	mu sync.RWMutex
	m  map[string]any
}

type ShardedMap []*Shard

func NewShardedMap(n int) ShardedMap {
	shards := make([]*Shard, n)
	for i := range shards {
		shards[i] = &Shard{
			m: make(map[string]any),
		}
	}

	return shards
}

func (m ShardedMap) Get(key string) any {
	shard := m.getShard(key)

	shard.mu.RLock()
	defer shard.mu.RUnlock()

	return shard.m[key]
}

func (m ShardedMap) Set(key string, value any) {
	shard := m.getShard(key)

	shard.mu.Lock()
	defer shard.mu.Unlock()

	shard.m[key] = value
}

func (m ShardedMap) Delete(key string) {
	shard := m.getShard(key)

	shard.mu.Lock()
	defer shard.mu.Unlock()

	delete(shard.m, key)
}

func (m ShardedMap) Keys() []string {
	var wg sync.WaitGroup
	var keys []string

	keysCh := make(chan string)
	wg.Add(len(m))

	for _, shard := range m {
		go func(s *Shard) {
			s.mu.RLock()
			for key := range s.m {
				keysCh <- key
			}
			s.mu.RUnlock()

			wg.Done()
		}(shard)
	}

	go func() {
		wg.Wait()
		close(keysCh)
	}()

	for key := range keysCh {
		keys = append(keys, key)
	}

	return keys
}

func (m ShardedMap) getShard(key string) *Shard {
	checksum := sha1.Sum([]byte(key))
	hash := int(checksum[17])
	index := hash % len(m)

	return m[index]
}
