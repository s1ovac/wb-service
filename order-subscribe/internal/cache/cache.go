package cache

import (
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/s1ovac/order-subscribe/internal/store/databases/order"
)

type Cache struct {
	model *order.Order
	cache *cache.Cache
}

func NewCache(model *order.Order) *Cache {
	return &Cache{
		model: model,
		cache: cache.New(5*time.Minute, 10*time.Minute),
	}
}

func (c *Cache) InitCache(orderID string) error {
	err := c.cache.Add(orderID, c.model, cache.DefaultExpiration)
	if err != nil {
		return err
	}
	return nil
}

func (c *Cache) Cache(k string) *order.Order {
	ord, found := c.cache.Get(k)
	if !found {
		return nil
	}
	return ord.(*order.Order)
}
