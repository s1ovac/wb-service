package cache

import (
	"context"
	"errors"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/s1ovac/order-subscribe/internal/store/databases/order"
)

type Cache struct {
	repository order.Repository
	cache      *cache.Cache
}

func NewCache(rep order.Repository) *Cache {
	return &Cache{
		repository: rep,
		cache:      cache.New(5*time.Minute, 10*time.Minute),
	}
}

func (c *Cache) InitCache(ctx context.Context) error {
	orders, err := c.repository.FindAll(ctx)
	if err != nil {
		return err
	}
	for _, order := range orders {
		_, found := c.cache.Get(order.OrderUID)
		if !found {
			if err := c.cache.Add(order.OrderUID, order, cache.DefaultExpiration); err != nil {
				return err
			}
		}

	}
	return nil
}

func (c *Cache) GetCache(ctx context.Context, k string) (order.Order, error) {
	foo, found := c.cache.Get(k)
	if !found {
		var o order.Order
		o, err := c.repository.FindOne(ctx, k)
		if err != nil {
			return order.Order{}, err
		}
		if err := c.cache.Add(k, o, cache.DefaultExpiration); err != nil {
			return order.Order{}, err
		}
		return o, err
	}
	o, ok := foo.(order.Order)
	if !ok {
		return order.Order{}, errors.New("undefined type of order")
	}
	return o, nil
}
