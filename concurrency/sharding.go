package concurrency

import (
	"crypto/sha1"
	"sync"
)

// Shard is an individually lockable collection representing a single data partition.
type Shard struct {
	mu sync.RWMutex
	m  map[string]any
}

// ShardedMap implements vertical-sharding pattern. It splits a large data structure into
// multiple partition to localize the effects of read/write locks.
//
// It's an abstraction around one or more Shard providing read and write access
// if there was a single map. Whenever a value is read or written to the map abstraction,
// a hash value is calculated for the key and taking into account the number of shards determines
// the corresponding shard index. This allows to isolate the necessary locking to only the shard
// at that index.
type ShardedMap []*Shard

// NewShardedMap constructs ShardedMap. It takes the number of shards among which
// the keys will be distributed.
func NewShardedMap(n int) ShardedMap {
	shards := make([]*Shard, n)
	for i := range shards {
		shards[i] = &Shard{
			m: make(map[string]any),
		}
	}

	return shards
}

// Get gets value by key.
func (m ShardedMap) Get(key string) any {
	shard := m.getShard(key)

	shard.mu.RLock()
	defer shard.mu.RUnlock()

	return shard.m[key]
}

// Set sets value by key.
func (m ShardedMap) Set(key string, value any) {
	shard := m.getShard(key)

	shard.mu.Lock()
	defer shard.mu.Unlock()

	shard.m[key] = value
}

// Delete deletes value buy key.
func (m ShardedMap) Delete(key string) {
	shard := m.getShard(key)

	shard.mu.Lock()
	defer shard.mu.Unlock()

	delete(shard.m, key)
}

// Keys returns all the existed keys.
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
