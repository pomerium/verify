package verify

import (
	lru "github.com/hashicorp/golang-lru"
	"github.com/rs/zerolog/log"

	"github.com/pomerium/sdk-go"
)

var _ sdk.JSONWebKeyStore = &Cache{}

// Cache is a naive implementation of a JSONWebKeyStore using a LRU cache.
type Cache struct{ lru *lru.Cache }

// NewCache creates a new cache.
func NewCache(size int) *Cache {
	c, err := lru.New(size)
	if err != nil {
		// only errors on invalid size
		log.Fatal().Err(err).Send()
	}
	return &Cache{c}
}

// Get gets a value from the cache.
func (c *Cache) Get(key interface{}) (value interface{}, ok bool) {
	return c.lru.Get(key)
}

// Add adds a value to the cache.
func (c *Cache) Add(key, value interface{}) {
	_ = c.lru.Add(key, value) // we don't care about eviction
}
