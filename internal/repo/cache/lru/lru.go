package lru

import (
	"container/list"
	"sync"
	"time"

	"github.com/andreyxaxa/order_svc/internal/entity"
)

type item struct {
	key      string
	value    entity.Order
	expireAt time.Time
}

type LRUCache struct {
	capacity int
	ttl      time.Duration
	mu       sync.Mutex
	items    map[string]*list.Element
	queue    *list.List
}

func New(capacity int, ttl time.Duration) *LRUCache {
	return &LRUCache{
		capacity: capacity,
		ttl:      ttl,
		items:    make(map[string]*list.Element),
		queue:    list.New(),
	}
}

func (c *LRUCache) Set(key string, value entity.Order) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Если уже есть - просто переместим в начало.
	if elem, ok := c.items[key]; ok {
		c.queue.MoveToFront(elem)
		return
	}

	it := &item{
		key:      key,
		value:    value,
		expireAt: time.Now().Add(c.ttl),
	}
	elem := c.queue.PushFront(it)
	c.items[key] = elem

	if c.queue.Len() > c.capacity {
		last := c.queue.Back()
		if last != nil {
			c.removeElem(last)
		}
	}
}

func (c *LRUCache) Get(key string) (entity.Order, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, ok := c.items[key]; ok {
		it := elem.Value.(*item)
		if time.Now().After(it.expireAt) {
			c.removeElem(elem)
			return entity.Order{}, false
		}
		c.queue.MoveToFront(elem)
		return it.value, true
	}

	return entity.Order{}, false
}

func (c *LRUCache) removeElem(el *list.Element) {
	it := el.Value.(*item)
	delete(c.items, it.key)
	c.queue.Remove(el)
}
