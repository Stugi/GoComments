package cache

import (
	"time"
)

type Cache interface {
	Get(key string) (any, bool)
	Set(key string, value any, ttl time.Duration) error
	Delete(key string) error
}

type Impl struct {
	cache map[string]any
	ttl   time.Duration
}

func New(ttl time.Duration) *Impl {
	return &Impl{
		cache: make(map[string]any),
		ttl:   ttl,
	}
}

func (c *Impl) Get(key string) (any, bool) {
	value, ok := c.cache[key]
	return value, ok
}
func (c *Impl) Set(key string, value any, ttl time.Duration) error {
	c.cache[key] = value
	go func() {
		time.Sleep(ttl)
		delete(c.cache, key)
	}()
	return nil
}
func (c *Impl) Delete(key string) error {
	delete(c.cache, key)
	return nil
}
