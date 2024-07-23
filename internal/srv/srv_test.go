package srv

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"lru-cache/internal/cache"
	"lru-cache/internal/models"
	"lru-cache/pkg/levelhandler"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func TestPutKeyHandler(t *testing.T) {
	mockCache := cache.NewMockCache()
		logger := slog.New(slog.NewTextHandler(os.Stdout, nil))


	server := &Server{
		storage: mockCache,
		cfg: Config{
			DefaultTTL: 10 * time.Minute,
		},
		logger: logger,
	}

	reqBody := `{"key":"testKey","value":"testValue","ttl_seconds":3600}`
	req := httptest.NewRequest(http.MethodPut, "/api/lru", strings.NewReader(reqBody))
	rec := httptest.NewRecorder()

	handler := http.HandlerFunc(server.postKey)
	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Equal(t, "testValue", mockCache.Store["testKey"])
}


func TestGetKeyHandler(t *testing.T) {
	mockCache := cache.NewMockCache()
	mockCache.Put(context.Background(), "testKey", "testValue", 3600*time.Second)

	levelHandler, err := levelhandler.New("DEBUG", slog.NewTextHandler(os.Stdout, nil))
	assert.NoError(t, err)

	logger := slog.New(levelHandler)

	server := &Server{
		storage: mockCache,
		logger:  logger,
	}

	router := chi.NewRouter()
	router.Get("/api/lru/{key}", server.getKey)

	req := httptest.NewRequest(http.MethodGet, "/api/lru/testKey", nil)
	rec := httptest.NewRecorder()

	routeCtx := chi.NewRouteContext()
	routeCtx.URLParams.Add("key", "testKey")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp models.GetResponse
	err = resp.FromJSON(rec.Body)
	assert.NoError(t, err)
	assert.Equal(t, "testKey", resp.Key)
	assert.Equal(t, "testValue", resp.Value)
}

func TestGetAllKeysHandler(t *testing.T) {
	mockCache := cache.NewMockCache()
	mockCache.Put(context.Background(), "testKey1", "testValue1", 3600*time.Second)
	mockCache.Put(context.Background(), "testKey2", "testValue2", 3600*time.Second)

	levelhandler, err := levelhandler.New("DEBUG", slog.NewTextHandler(os.Stdout, nil))

	assert.NoError(t,err)
	logger := slog.New(levelhandler)

	server := &Server{
		storage: mockCache,
		logger: logger,
	}

	req := httptest.NewRequest(http.MethodGet, "/keys", nil)
	rec := httptest.NewRecorder()

	handler := http.HandlerFunc(server.getAllKeys)
	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp models.GetAllResponse
	err = json.NewDecoder(rec.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.ElementsMatch(t, []string{"testKey1", "testKey2"}, resp.Keys)
	assert.ElementsMatch(t, []interface{}{"testValue1", "testValue2"}, resp.Values)	
}

func TestEvictKeyHandler(t *testing.T) {
	mockCache := cache.NewMockCache()
	mockCache.Put(context.Background(), "testKey", "testValue", 3600*time.Second)
	levelhandler, err := levelhandler.New("DEBUG", slog.NewTextHandler(os.Stdout, nil))

	assert.NoError(t,err)
	logger := slog.New(levelhandler)

	server := &Server{
		storage: mockCache,
		logger:  logger,
	}

	router := chi.NewRouter()
	router.Delete("/api/lru/{key}", server.evictKey)

	req := httptest.NewRequest(http.MethodDelete, "/api/lru/testKey", nil)
	rec := httptest.NewRecorder()

	routeCtx := chi.NewRouteContext()
	routeCtx.URLParams.Add("key", "testKey")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNoContent, rec.Code)
	assert.NotContains(t, mockCache.Store, "testKey")
}

func TestEvictAllKeysHandler(t *testing.T) {
	mockCache := cache.NewMockCache()
	mockCache.Put(context.Background(), "testKey1", "testValue1", 3600*time.Second)
	mockCache.Put(context.Background(), "testKey2", "testValue2", 3600*time.Second)
	levelhandler, err := levelhandler.New("DEBUG", slog.NewTextHandler(os.Stdout, nil))

	assert.NoError(t,err)
	logger := slog.New(levelhandler)
	server := &Server{
		storage: mockCache,
		logger: logger,
	}

	req := httptest.NewRequest(http.MethodDelete, "/api/lru", nil)
	rec := httptest.NewRecorder()

	handler := http.HandlerFunc(server.evictAllKeys)
	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNoContent, rec.Code)
	assert.Empty(t, mockCache.Store)
}
