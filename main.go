package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/redis/go-redis/v9"
)

var (
	rdb *redis.Client
)

func handleHealth(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "OK")
}

func redisConnect() {
	for {
		log.Println("Connection to Redis")
		rdb = redis.NewClient(&redis.Options{
			Addr:     "redis:6379",
			Password: "",
			DB:       0,
		})

		ctx := context.Background()
		_, err := rdb.Ping(ctx).Result()
		if err != nil {
			log.Printf("Not connected %v\n", err)
			time.Sleep(5 * time.Second)
		} else {
			break
		}
	}
}

func main() {
	fmt.Println("Starting")
	redisConnect()

	r := mux.NewRouter()
	r.HandleFunc("/health", handleHealth)
	http.Handle("/", r)
	err := http.ListenAndServe(":8080", r)
	log.Fatal(err)
}
