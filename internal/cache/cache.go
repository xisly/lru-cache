// Package cache provides an implementation of an LRU cache with concurrency support.
// The cache supports only string keys and simple data types as values.
package cache

import (
	"context"
	"lru-cache/pkg/errs"
	"sync"
	"time"
)

// ILRUCache интерфейс LRU-кэша. Поддерживает только строковые ключи. Поддерживает только простые типы данных в значениях.
type ILRUCache interface {
	// Put запись данных в кэш
	Put(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	// Get получение данных из кэша по ключу
	Get(ctx context.Context, key string) (value interface{}, expiresAt time.Time, err error)
	// GetAll получение всего наполнения кэша в виде двух слайсов: слайса ключей и слайса значений. Пары ключ-значения из кэша располагаются на соответствующих позициях в слайсах.
	GetAll(ctx context.Context) (keys []string, values []interface{}, err error)
	// Evict ручное удаление данных по ключу
	Evict(ctx context.Context, key string) (value interface{}, err error)
	// EvictAll ручная инвалидация всего кэша
	EvictAll(ctx context.Context) error
}

type node struct {
	key       string
	value     interface{}
	expiresAt time.Time
	prev      *node
	next      *node
}

type cache struct {
	capacity int
	data     map[string]*node
	left     *node
	right    *node
  mu sync.Mutex
}

// New creates a new LRU cache with the specified capacity.
// Returns an instance of the ILRUCache interface.
func New(capacity int) ILRUCache {
	ret := cache{
		capacity: capacity,
		data:     make(map[string]*node, capacity),
		left:     &node{key: "", value: ""},
		right:    &node{key: "", value: ""},
	}
	ret.left.next, ret.right.prev = ret.right, ret.left
	return &ret
}

// Put stores data in the cache with a specified TTL.
// If the key already exists, the existing entry is updated.
func (c *cache) Put(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
  c.mu.Lock()
  defer c.mu.Unlock()

	if _, ok := c.data[key]; ok {
		c.remove(c.data[key])
	}
	c.data[key] = &node{key: key, value: value, expiresAt: time.Now().Add(ttl)}
	c.insert(c.data[key])

	if len(c.data) > c.capacity {
		nd := c.left.next
		c.remove(nd)
		delete(c.data, nd.key)
	}
	return nil
}

func (c *cache) insert(nd *node) {
	prev, nxt := c.right.prev, c.right
	nxt.prev = nd
	prev.next = nxt.prev
	nd.next, nd.prev = nxt, prev
}

func (c *cache) remove(nd *node) {
	prev, nxt := nd.prev, nd.next
	prev.next, nxt.prev = nxt, prev
}

// Get retrieves data from the cache by key.
// Returns the value, expiration time, and an error if the key is not found or has expired.
func (c *cache) Get(ctx context.Context, key string) (value interface{}, expiresAt time.Time, err error) {
  c.mu.Lock()
  defer c.mu.Unlock()

	if nd, ok := c.data[key]; ok {
		if nd.expiresAt.Before(time.Now()) {
			c.remove(nd)
			delete(c.data, key)
			return nil, time.Time{}, errs.ErrNotFound
		}
		c.remove(nd)
		c.insert(nd)
		return nd.value, nd.expiresAt, nil
	}
	return nil, time.Time{}, errs.ErrNotFound
}

// GetAll retrieves all entries from the cache as two slices: a slice of keys and a slice of values.
// Returns an error if the cache is empty.
func (c *cache) GetAll(ctx context.Context) (keys []string, values []interface{}, err error) {
  c.mu.Lock()
  defer c.mu.Unlock()

	if len(c.data) == 0 {
		err = errs.ErrCacheIsEmpty
		return
	}
	keys = make([]string, len(c.data))
	values = make([]interface{}, len(c.data))
	i := 0
	for k, v := range c.data {
		keys[i], values[i] = k, v.value
		i++
	}
	return keys, values, nil
}

// Evict manually removes data by key from the cache.
// Returns the value and an error if the key is not found or has expired.
func (c *cache) Evict(ctx context.Context, key string) (value interface{}, err error) {
  c.mu.Lock()
  defer c.mu.Unlock()

	if len(c.data) == 0 {
		err = errs.ErrCacheIsEmpty
		return
	}
	if nd, ok := c.data[key]; ok {
		if nd.expiresAt.Before(time.Now()) {
			c.remove(nd)
			delete(c.data, key)
			return nil, errs.ErrNotFound
		}
		c.remove(nd)
		delete(c.data, nd.key)
		return nd.value, nil
	}

	return nil, errs.ErrNotFound
}

// EvictAll manually invalidates the entire cache.
func (c *cache) EvictAll(ctx context.Context) error {
  c.mu.Lock()
  defer c.mu.Unlock()

	c.left.next, c.right.prev = c.right, c.left
	clear(c.data)

	return nil
}