package lru

import "errors"

var (
	cache = make(map[string]*Item)
)

func Set(item *Item) {
	cache[item.Key] = item
}

func Get(key string) (*Item, error) {
	if value, ok := cache[key]; ok {
		return value, nil
	}

	return nil, errors.New("not found")
}

func Delete(key string) {
	delete(cache, key)
}

func Flush() {
	cache = make(map[string]*Item)
}
