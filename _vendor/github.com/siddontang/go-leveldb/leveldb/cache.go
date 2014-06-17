package leveldb

import (
	"github.com/siddontang/goleveldb/leveldb/cache"
)

type Cache struct {
	c cache.Cache
}

func NewLRUCache(capacity int) *Cache {
	return &Cache{cache.NewLRUCache(capacity)}
}

func (c *Cache) Close() {

}
