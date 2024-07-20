package cache

import (
	"context"
	"fmt"
	"lru-cache/pkg/errs"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

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

func TestGet(t *testing.T) {
	cache := New(23248247)

	for i := range 5 {
		err := cache.Put(context.Background(), fmt.Sprintf("%d", i), i, time.Minute)
		assert.NoError(t, err)
	}

	for i := range 5 {
		gotValue, _, err := cache.Get(context.Background(), fmt.Sprintf("%d", i))
		assert.NoError(t, err)
		assert.Equal(t, i, gotValue)
	}
}

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
