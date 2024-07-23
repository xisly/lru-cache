package cache

import (
	"context"
	"fmt"
	"lru-cache/pkg/errs"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestLazyTTL verifies that items with TTL are correctly removed from the cache after the TTL expires.
func TestLazyTTL(t *testing.T) {
	cache := New(5)

	for i := range 5 {
		err := cache.Put(context.Background(), fmt.Sprintf("%d", i), i, time.Millisecond*15)
		assert.NoError(t, err)
	}

	time.Sleep(time.Millisecond * 100)

	for i := range 5 {
		got, _, err := cache.Get(context.Background(), fmt.Sprintf("%d", i))
		assert.Equal(t, errs.ErrNotFound, err)
		assert.Nil(t, got)
	}

}

// TestPut verifies that items can be added to the cache and that an existing item can be updated.
func TestPut(t *testing.T) {
	cache := New(5)

	for i := range 5 {
		err := cache.Put(context.Background(), fmt.Sprintf("%d", i), i, time.Hour)
		assert.NoError(t, err)
	}

	prevValue, prevExpirationTime, err := cache.Get(context.Background(), "0")
	assert.NoError(t, err)

	cache.Put(context.Background(), "0", "hello world", time.Hour*2)
	actualValue, actualExpirationTime, err := cache.Get(context.Background(), "0")

	assert.NoError(t, err)
	assert.NotEqual(t, prevValue, actualValue)
	assert.NotEqual(t, prevExpirationTime, actualExpirationTime)
}

// TestEvict verifies that an item can be manually removed from the cache and that the item is no longer retrievable.
func TestEvict(t *testing.T) {
	cache := New(5)

	for i := range 5 {
		cache.Put(context.Background(), fmt.Sprintf("%d", i), i, time.Hour)
	}

	want := 3
	got, err := cache.Evict(context.Background(), "3")

	assert.NoError(t, err)
	assert.Equal(t, want, got)

	_, _, err = cache.Get(context.Background(), "3")
	assert.Equal(t, errs.ErrNotFound, err)

}

// TestEvictAll verifies that all items can be removed from the cache and that the cache is empty afterwards.
func TestEvictAll(t *testing.T) {

	cache := New(5)
	for i := range 5 {
		cache.Put(context.Background(), fmt.Sprintf("%d", i), i, time.Hour)
	}

	err := cache.EvictAll(context.Background())
	assert.NoError(t, err)

	var (
		wantKeys   []string
		wantValues []interface{}
	)
	gotKeys, gotValues, err := cache.GetAll(context.Background())
	assert.Equal(t, errs.ErrCacheIsEmpty, err)
	assert.Equal(t, wantKeys, gotKeys)
	assert.Equal(t, wantValues, gotValues)
}

// TestGet verifies that an item can be retrieved from the cache and that the correct value and expiration time are returned.
func TestGet(t *testing.T) {
	cache := New(23248247)

	for i := range 5 {
		err := cache.Put(context.Background(), fmt.Sprintf("%d", i), i, time.Minute)
		assert.NoError(t, err)
	}

	for i := range 5 {
		gotValue, gotTime, err := cache.Get(context.Background(), fmt.Sprintf("%d", i))
		assert.NoError(t, err)
		assert.Equal(t, i, gotValue)
		assert.NotEqual(t, time.Time{}, gotTime)
	}
}

// TestGetAll verifies that all items can be retrieved from the cache and that the correct keys and values are returned.
func TestGetAll(t *testing.T) {
	wantKeys := []string{"0", "1", "2", "3", "4"}
	wantValues := []interface{}{0, 1, 2, 3, 4}
	cache := New(5)
	for i := range 5 {
		cache.Put(context.Background(), fmt.Sprintf("%d", i), i, time.Hour)
	}
	gotKeys, gotValues, err := cache.GetAll(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, wantKeys, gotKeys)
	assert.Equal(t, wantValues, gotValues)

}

// TestLeastRecentlyAddedKeyIsEvicted verifies that the least recently added key is evicted when the cache exceeds its capacity.
func TestLeastRecntlyAddedKeyIsEvicted(t *testing.T) {
	cache := New(5)
	for i := range 6 {
		cache.Put(context.Background(), fmt.Sprintf("%d", i), i, time.Hour)
	}
	value, expiresAt, err := cache.Get(context.Background(), "0")
	assert.Equal(t, nil, value)
	assert.Equal(t, time.Time{}, expiresAt)
	assert.Equal(t, errs.ErrNotFound, err)
}
