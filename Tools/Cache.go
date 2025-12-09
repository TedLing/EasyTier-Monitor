package Tools

import (
	"sync"
	"time"
)

// CacheItem 缓存项结构体
type CacheItem struct {
	Value      interface{}
	Expiration int64
	CreatedAt  time.Time
}

// Cache 缓存结构体
type Cache struct {
	items map[string]CacheItem
	mutex sync.RWMutex
}

// NewCache 创建一个新的缓存实例
func NewCache() *Cache {
	return &Cache{
		items: make(map[string]CacheItem),
	}
}

// Set 设置缓存项，duration为缓存过期时间
func (c *Cache) Set(key string, value interface{}, duration time.Duration) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	var expiration int64
	if duration > 0 {
		expiration = time.Now().Add(duration).UnixNano()
	} else {
		// 永不过期
		expiration = 0
	}

	c.items[key] = CacheItem{
		Value:      value,
		Expiration: expiration,
		CreatedAt:  time.Now(),
	}
}

// Get 获取缓存项，如果不存在或已过期则返回nil
func (c *Cache) Get(key string) (interface{}, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	item, found := c.items[key]
	if !found {
		return nil, false
	}

	// 检查是否过期
	if item.Expiration > 0 && time.Now().UnixNano() > item.Expiration {
		return nil, false
	}

	return item.Value, true
}

// Delete 删除缓存项
func (c *Cache) Delete(key string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	delete(c.items, key)
}

// Clear 清空所有缓存项
func (c *Cache) Clear() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.items = make(map[string]CacheItem)
}

// Size 获取当前缓存项数量
func (c *Cache) Size() int {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return len(c.items)
}

// 全局缓存实例
var AppCache = NewCache()

// 缓存键常量
const (
	CacheKeyPeer      = "peer_info"
	CacheKeyNode      = "node_info"
	CacheKeyConnector = "connector_info"
)

// 默认缓存过期时间
var (
	DefaultCacheDuration   = 10 * time.Second // 默认缓存10秒
	PeerCacheDuration      = 10 * time.Second
	NodeCacheDuration      = 30 * time.Second
	ConnectorCacheDuration = 30 * time.Second
)

// GetCacheDuration 根据配置获取缓存过期时间
func GetCacheDuration(key string) time.Duration {
	// 如果缓存功能被禁用，直接返回0（不缓存）
	if !AppConfig.Cache.Enabled {
		return 0
	}

	// 根据不同的缓存键返回不同的过期时间
	switch key {
	case CacheKeyPeer:
		if AppConfig.Cache.PeerTTL > 0 {
			return time.Duration(AppConfig.Cache.PeerTTL) * time.Second
		}
		return PeerCacheDuration
	case CacheKeyNode:
		if AppConfig.Cache.NodeTTL > 0 {
			return time.Duration(AppConfig.Cache.NodeTTL) * time.Second
		}
		return NodeCacheDuration
	case CacheKeyConnector:
		if AppConfig.Cache.ConnectorTTL > 0 {
			return time.Duration(AppConfig.Cache.ConnectorTTL) * time.Second
		}
		return ConnectorCacheDuration
	default:
		if AppConfig.Cache.DefaultTTL > 0 {
			return time.Duration(AppConfig.Cache.DefaultTTL) * time.Second
		}
		return DefaultCacheDuration
	}
}
