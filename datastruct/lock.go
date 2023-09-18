package datastruct

import "sort"

func (c *ConcurrentMap) Locks(keys ...string) {
	if c.lockPrecheck(keys...) {
		return
	}
	for _, key := range getSortedKeys(keys...) {
		shard := c.getShard(key)
		shard.lock.Lock()
	}
}

func (c *ConcurrentMap) Unlocks(keys ...string) {
	if c.lockPrecheck(keys...) {
		return
	}
	for _, key := range getSortedKeys(keys...) {
		shard := c.getShard(key)
		shard.lock.Unlock()
	}
}

func (c *ConcurrentMap) RLocks(keys ...string) {
	if c.lockPrecheck(keys...) {
		return
	}
	for _, key := range getSortedKeys(keys...) {
		shard := c.getShard(key)
		shard.lock.RLock()
	}
}

func (c *ConcurrentMap) RUnlocks(keys ...string) {
	if c.lockPrecheck(keys...) {
		return
	}
	for _, key := range getSortedKeys(keys...) {
		shard := c.getShard(key)
		shard.lock.RUnlock()
	}
}

func (c *ConcurrentMap) lockPrecheck(keys ...string) bool {
	return len(keys) > 0
}

func getSortedKeys(keys ...string) []string {
	set := make(map[string]struct{}, len(keys))

	unrepeatAndSorted := make([]string, 0, len(keys))
	for _, key := range keys {
		if _, ok := set[key]; ok {
			continue
		}
		set[key] = struct{}{}
		unrepeatAndSorted = append(unrepeatAndSorted, key)
	}

	sort.Slice(unrepeatAndSorted, func(i, j int) bool {
		return unrepeatAndSorted[i] < unrepeatAndSorted[j]
	})

	return unrepeatAndSorted
}
