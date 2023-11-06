package main

import (
	"context"
	"husir_blog/src/api"
	"husir_blog/src/db"
	"husir_blog/src/rdb"
	"husir_blog/src/util"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {

	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("can not load config: ", err)
	}

	clientOpt := options.Client().ApplyURI(config.DBSource)
	client, err := mongo.Connect(context.Background(), clientOpt)
	if err != nil {
		log.Fatal("can not connect to db")
	}
	database := client.Database("blog")
	store := db.NewStore(database)

	// rdb := redis.NewClient(&redis.Options{
	// 	Addr:     "localhost:6379",
	// 	Password: "", // no password set
	// 	DB:       0,  // use default DB
	// })

	rdb, err := rdb.NewRedisClient()
	if err != nil {
		log.Fatal("can not connect redis: ", err)

	}

	log.Println(rdb)

	server, err := api.NewServer(config, store, rdb)
	if err != nil {
		log.Fatal("can not create sever: ", err)

	}

	err = server.CachePostRecommend(context.Background())
	if err != nil {
		log.Printf(">> CachePostRecommend failed: %s", err)

	} else {
		log.Printf(">> CachePostRecommend success")

	}

	err = server.Start(config.HttpServerAddress)
	if err != nil {
		log.Fatal("can not start server: ", err)
	}

	// go func(s *api.Server) {
	// 	err := s.CachePostRecommend(context.Background())
	// 	if err != nil {
	// 		log.Printf(">> CachePostRecommend failed: %s", err)

	// 	} else {
	// 		log.Printf(">> CachePostRecommend success")

	// 	}
	// }(server)

	// time.Sleep(time.Second)

}
