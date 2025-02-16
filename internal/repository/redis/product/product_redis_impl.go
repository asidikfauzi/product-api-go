package product

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"product-api-go/internal/handler/product/dto"
	"product-api-go/internal/pkg/constant"
)

var ctx = context.Background()

type products struct {
	redis *redis.Client
}

func NewProductsRedis(r *redis.Client) ProductsRedis {
	return &products{
		redis: r,
	}
}

func (c *products) GetAllProduct(key string) (res dto.ProductsResponseWithPage, err error) {

	value, err := c.redis.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return res, constant.KeyRedisNotExists
	} else if err != nil {
		return res, err
	}

	err = json.Unmarshal([]byte(value), &res)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (c *products) CreateAllProduct(key string, data dto.ProductsResponseWithPage) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	err = c.redis.Set(ctx, key, jsonData, 0).Err()
	if err != nil {
		return err
	}

	return nil
}

func (c *products) GetProductById(uuid uuid.UUID) (res dto.ProductResponse, err error) {
	key := fmt.Sprintf(constant.ProductByIdKey, uuid.String())

	value, err := c.redis.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return res, constant.KeyRedisNotExists
	} else if err != nil {
		return res, err
	}

	err = json.Unmarshal([]byte(value), &res)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (c *products) CreateProductById(uuid uuid.UUID, data dto.ProductResponse) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	key := fmt.Sprintf(constant.ProductByIdKey, uuid.String())

	err = c.redis.Set(ctx, key, jsonData, 0).Err()
	if err != nil {
		return err
	}

	return nil
}

func (c *products) DeleteAll(key string) error {
	var cursor uint64
	var keys []string
	var err error

	for {
		var result []string
		result, cursor, err = c.redis.Scan(ctx, cursor, key+"*", 1000).Result()
		if err != nil {
			return err
		}

		keys = append(keys, result...)

		if cursor == 0 {
			break
		}
	}

	if len(keys) == 0 {
		return nil
	}

	_, err = c.redis.Del(ctx, keys...).Result()
	if err != nil {
		return err
	}

	return nil
}

func (c *products) Delete(key string) error {
	err := c.redis.Del(ctx, key).Err()
	if err != nil {
		return err
	}

	return nil
}
