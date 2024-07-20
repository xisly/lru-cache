package cache 

import (
  "sync"
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
  key string 
  value interface{}
  expiresAt time.Duration
  prev *node
  next *node
}

type cache struct {
  capacity int
  cache map[string]*node
  left *node
  right *node
  mu sync.Mutex
}

func New(capacity int) ILRUCache {
  ret := cache{
    capacity: capacity
    cache: make(map[string]node,capacity)
    left: &node{key:"", value:""}
    right: &node{key:"", value: ""}
  }
  ret.left.next, ret.right.prev = ret.right, ret.left
  return ret
}

func (c *cache) Put(ctx context.Context, key string, value interface{} ttl time.Duration) error{
  return nil
}

func (c *cache) insert(nd *node) {
  prev, nxt := this.right.prev, this.right
  nxt.prev = nd
  prev.next = nxt.prev
  node.next, node.prev = nxt, prev
}

func (this *LRUCache) remove(nd *node) {
    prev, nxt := nd.prev, nd.next
    prev.next, nxt.prev = nxt, prev
}

func (c *cache) Get(ctx context.Context, key string) (value interface{}, expiresAt time.Time, err error) {
   if _, ok := this.cache[key]; ok {
        this.Remove(this.cache[key])
        this.Insert(this.cache[key])
        return this.cache[key].val
    }
    return nil, time.Now(), errors.New()
}

func (c *cache) GetAll(ctx context.Context) (keys []string, values []interface{}, err error)

func (c *cache) Evict(ctx context.Context, key string) (value interface{}, err error)

func (c *cache) EvictAll(ctx context.Context) error {} 

