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

func handleHome(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Key Value Store v1.0.0")
}

func handleAdd(w http.ResponseWriter, r *http.Request) {
	if r.Method == "PUT" {
		vars := mux.Vars(r)
		key := vars["key"]
		val := vars["value"]

		ctx := context.Background()
		_, err := rdb.Set(ctx, key, val, 0).Result()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Error %v\n", err)
		}

		fmt.Fprintln(w, "OK")
	} else {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Wrong method")
	}
}

func handleList(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		ctx := context.Background()
		keys, err := rdb.Keys(ctx, "*").Result()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error %v\n", err)
		} else {
			for _, key := range keys {
				fmt.Fprintln(w, key)
			}
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Wrong method")
	}
}

func handleFetch(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		vars := mux.Vars(r)
		key := vars["key"]

		ctx := context.Background()
		val, err := rdb.Get(ctx, key).Result()
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintln(w, "Not found")
		} else {
			fmt.Fprintln(w, val)
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Wrong method")
	}
}

func handleDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method == "DEL" {
		vars := mux.Vars(r)
		key := vars["key"]

		ctx := context.Background()
		_, err := rdb.Del(ctx, key).Result()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error %v\n", err)
		} else {
			fmt.Fprintln(w, "OK")
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Wrong method")
	}
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
			log.Println("Connected to Redis")
			break
		}
	}
}

func main() {
	fmt.Println("Starting")
	redisConnect()

	r := mux.NewRouter()
	r.HandleFunc("/", handleHome)
	r.HandleFunc("/health", handleHealth)
	r.HandleFunc("/add/{key}/{value}", handleAdd)
	r.HandleFunc("/fetch/{key}", handleFetch)
	r.HandleFunc("/list", handleList)
	r.HandleFunc("/del/{key}", handleDelete)
	http.Handle("/", r)
	err := http.ListenAndServe(":8080", r)
	log.Fatal(err)
}
