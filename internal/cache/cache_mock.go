package cache

import (
	"context"
	"lru-cache/pkg/errs"
	"time"
)

// MockCache is a custom mock implementation of the ILRUCache interface.
// It stores key-value pairs and their corresponding TTLs.
type MockCache struct {
	Store    map[string]interface{}
	TTLStore map[string]time.Time
}

// NewMockCache creates a new instance of MockCache.
// Returns a pointer to the newly created MockCache.
func NewMockCache() *MockCache {
	return &MockCache{
		Store:    make(map[string]interface{}),
		TTLStore: make(map[string]time.Time),
	}
}

// Put stores a key-value pair in the cache with a specified TTL.
// If the TTL is greater than zero, the expiration time is calculated and stored.
func (m *MockCache) Put(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	m.Store[key] = value
	if ttl > 0 {
		m.TTLStore[key] = time.Now().Add(ttl)
	}
	return nil
}

// Get retrieves a value from the cache by key.
// Returns the value, its expiration time, and an error if the key is not found.
func (m *MockCache) Get(ctx context.Context, key string) (interface{}, time.Time, error) {
	if value, ok := m.Store[key]; ok {
		expiry := m.TTLStore[key]
		return value, expiry, nil
	}
	return nil, time.Time{}, errs.ErrNotFound
}

// GetAll retrieves all key-value pairs from the cache.
// Returns slices of keys and values, and an error if the cache is empty.
func (m *MockCache) GetAll(ctx context.Context) ([]string, []interface{}, error) {
	if len(m.Store) == 0 {
		return nil, nil, errs.ErrCacheIsEmpty
	}

	keys := make([]string, 0, len(m.Store))
	values := make([]interface{}, 0, len(m.Store))
	for k, v := range m.Store {
		keys = append(keys, k)
		values = append(values, v)
	}
	return keys, values, nil
}

// Evict removes a key-value pair from the cache by key.
// Returns the value and an error if the key is not found.
func (m *MockCache) Evict(ctx context.Context, key string) (interface{},error) {
	if _, ok := m.Store[key]; !ok {
		return nil,errs.ErrNotFound
	}
	v := m.Store[key]
	delete(m.Store, key)
	delete(m.TTLStore, key)
	return v, nil
}

// EvictAll removes all key-value pairs from the cache.
// Returns an error if the cache is already empty.
func (m *MockCache) EvictAll(ctx context.Context) error {
	if len(m.Store) == 0 {
		return errs.ErrCacheIsEmpty
	}
	m.Store = make(map[string]interface{})
	m.TTLStore = make(map[string]time.Time)
	return nil
}