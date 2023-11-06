package rdb

import (
	"context"
	"encoding/json"
	"husir_blog/src/db"
	"log"
	"sync"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	Client *redis.Client
	Mutex  sync.Mutex
}

func NewRedisClient() (*RedisClient, error) {
	client := RedisClient{}

	client.Mutex.Lock()
	defer client.Mutex.Unlock()

	if client.Client == nil {
		client.Client = redis.NewClient(&redis.Options{
			Addr:     "localhost:6377",
			Password: "",
		})
	}

	return &client, nil
}

func (client *RedisClient) GetClient() *redis.Client {
	return client.Client
}

func (client *RedisClient) DeleteKey(ctx context.Context, key string) {
	res, err := client.Client.Del(ctx, key).Result()
	if err != nil {
		log.Printf("DEL key %s failed: %s", key, err)
	}
	log.Printf("DEL key %s is %d", key, res)
}

func (client *RedisClient) GetListPost(ctx context.Context, key string, paging db.Paging) []*db.PostResponseFull {

	min := (paging.Page - 1) * paging.Limit
	max := min + paging.Limit - 1

	values, err := client.Client.LRange(ctx, key, int64(min), int64(max)).Result()

	if err != nil {
		log.Printf("> key:%s does not exists", key)
		return nil
	}

	if len(values) <= 0 {
		return nil
	}

	var list []*db.PostResponseFull
	for _, v := range values {
		var post *db.PostResponseFull
		err = json.Unmarshal([]byte(v), &post)
		if err != nil {
			log.Printf("> can not Unmarshal Json.  Err: %s", err)
			return nil
		}

		list = append(list, post)
	}

	log.Println(">> Post get from redis")

	return list

}

func (client *RedisClient) SetListPost(ctx context.Context, key string, data []*db.PostResponseFull) {

	for index, v := range data {
		post, err := json.Marshal(v)
		if err != nil {
			log.Printf("can not Marshal Json.  Err: %s", err)
		}
		result, err :=
			client.Client.LPush(ctx, key, string(post)).Result()
		if err != nil {
			log.Printf("> can not add post:%d  Err: %s", index, err)
		} else {
			log.Printf("> post %d added: %d", index, result)
		}
	}

}

func (client *RedisClient) SetPostsRecommend(ctx context.Context, key string, data []*db.PostRecommendResponse) {

	b, _ := json.Marshal(data)
	var m map[string]interface{}
	err := json.Unmarshal(b, &m)
	if err != nil {
		log.Println("errrrrr")
	}

	res, err := client.Client.HSet(ctx, key, m).Result()
	if err != nil {
		log.Println("errrrrr hset")
	}
	log.Print(res)
}
func (client *RedisClient) SetPostBySlug(ctx context.Context, post *db.PostResponseFull) {

	p, err := json.Marshal(post)
	if err != nil {
		log.Println(err)
	}
	res, err := client.Client.HSet(ctx, "post", []string{post.Slug, string(p)}).Result()
	if err != nil {
		log.Println(err)
	}
	log.Print(res)
}
func (client *RedisClient) GetPostBySlug(ctx context.Context, slug string) (*db.PostResponseFull, error) {
	data, err := client.Client.HGet(ctx, "post", slug).Result()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var p db.PostResponseFull

	err = json.Unmarshal([]byte(data), &p)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	log.Println(">> GetPostBySlug from redis")
	return &p, err
}
