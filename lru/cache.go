package lru

var (
	LRU *Cache
)

func CreateLRU(size int) {
	LRU = New(size)
}