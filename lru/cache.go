package lru

import (
	"errors"
	"runtime"
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
	delete(cache, key)
}

func Flush() {
	cache = make(map[string]*Item)

	runtime.GC()
}

func Has(key string) bool {
	_, ok := cache[key]
	return ok
}

func Clean() {
	for key, item := range cache {
		if item.IsExpired() {
			Delete(key)
		}
	}
}