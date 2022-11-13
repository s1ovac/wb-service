package cache

import (
	"context"
	"errors"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/s1ovac/order-subscribe/internal/store/databases/order"
	"github.com/sirupsen/logrus"
)

type Cache struct {
	repository order.Repository
	cache      *cache.Cache
	logger     *logrus.Logger
}

func NewCache(rep order.Repository, logger *logrus.Logger) *Cache {
	return &Cache{
		repository: rep,
		cache:      cache.New(5*time.Minute, 10*time.Minute),
		logger:     logger,
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
	c.logger.Info("getting order from cache")
	foo, found := c.cache.Get(k)
	if !found {
		c.logger.Info("order not found in cache")
		var o order.Order
		o, err := c.repository.FindOne(ctx, k)
		if err != nil {
			return order.Order{}, err
		}
		c.logger.Info("adding new order to the cache")
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

func (c *Cache) CheckCache(k string) bool {
	_, found := c.cache.Get(k)
	return found
}
