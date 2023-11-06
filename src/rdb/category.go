package rdb

import (
	"context"
	"encoding/json"
	"husir_blog/src/db"
	"log"
	"time"
)

func (c *RedisClient) SetCategories(ctx context.Context, key string, data []*db.CategoryResponse) {

	v, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error marshal category: %s", err)
	} else {
		res, err := c.Client.Set(ctx, key, string(v), time.Hour).Result()
		if err != nil {

			log.Printf("Can not set categories to redis: %s", err)

		}
		log.Printf("Add category to redis success: %s", res)
	}

}
func (c *RedisClient) GetCategories(ctx context.Context, key string) ([]*db.CategoryResponse, error) {

	res, err := c.Client.Get(ctx, key).Result()
	if err != nil {
		log.Printf("Can not get categories from redis: %s", err)
		return nil, err
	}

	var data []*db.CategoryResponse
	err = json.Unmarshal([]byte(res), &data)
	if err != nil {
		log.Printf("Can not unmarshal categories from redis: %s", err)
		return nil, err
	}

	log.Print("Get categories from redis success")
	return data, nil
}
