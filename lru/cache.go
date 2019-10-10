package lru

import (
	"errors"
	"sync"
)

var (
	cache = make(map[string]*Item)
	mutex sync.Mutex
)

func Set(item *Item) {
	cache[item.Key] = item
}

func Get(key string) (*Item, error) {
	value, ok := cache[key]
	if ok {
		return value, nil
	}

	return nil, errors.New("not found")
}

func Delete(key string) {
	mutex.Lock()
	delete(cache, key)
	mutex.Unlock()
}

func Flush() {
	mutex.Lock()
	cache = make(map[string]*Item)
	mutex.Unlock()
}
