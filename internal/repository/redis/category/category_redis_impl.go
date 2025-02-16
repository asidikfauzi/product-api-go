package category

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"log"
	"product-api-go/internal/handler/category/dto"
	"product-api-go/internal/pkg/constant"
)

var ctx = context.Background()

type categories struct {
	redis *redis.Client
}

func NewCategoriesRedis(r *redis.Client) CategoriesRedis {
	return &categories{
		redis: r,
	}
}

func (c *categories) GetAllCategory(key string) (res dto.CategoriesResponseWithPage, err error) {

	value, err := c.redis.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return res, errors.New(constant.KeyRedisNotExists)
	} else if err != nil {
		log.Printf(constant.RedisGetError, err)
		return res, err
	}

	err = json.Unmarshal([]byte(value), &res)
	if err != nil {
		log.Printf(constant.JsonUnmarshalError, err)
		return res, err
	}

	return res, nil
}

func (c *categories) CreateAllCategory(key string, data dto.CategoriesResponseWithPage) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Printf(constant.JsonMarshalError, err)
		return err
	}

	err = c.redis.Set(ctx, key, jsonData, 0).Err()
	if err != nil {
		log.Printf(constant.RedisSetError, err)
		return err
	}

	return nil
}

func (c *categories) GetCategoryById(uuid uuid.UUID) (res dto.CategoryResponse, err error) {
	key := fmt.Sprintf(constant.CategoryByIdKey, uuid.String())

	value, err := c.redis.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return res, errors.New(constant.KeyRedisNotExists)
	} else if err != nil {
		log.Printf(constant.RedisGetError, err)
		return res, err
	}

	err = json.Unmarshal([]byte(value), &res)
	if err != nil {
		log.Printf(constant.JsonUnmarshalError, err)
		return res, err
	}

	return res, nil
}

func (c *categories) CreateCategoryById(uuid uuid.UUID, data dto.CategoryResponse) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Printf(constant.JsonMarshalError, err)
		return err
	}

	key := fmt.Sprintf(constant.CategoryByIdKey, uuid.String())

	err = c.redis.Set(ctx, key, jsonData, 0).Err()
	if err != nil {
		log.Printf(constant.RedisSetError, err)
		return err
	}

	return nil
}

func (c *categories) DeleteAll(key string) error {
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

func (c *categories) Delete(key string) error {
	err := c.redis.Del(ctx, key).Err()
	if err != nil {
		log.Printf(constant.RedisDeleteError, err)
		return err
	}

	return nil
}
