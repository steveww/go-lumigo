package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
)

const (
	LUMIGO = "ga-otlp.lumigo-tracer-edge.golumigo.com"
)

var (
	rdb    *redis.Client
	tracer trace.Tracer
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

func appResource() *resource.Resource {
	return resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName("go-lumigo"),
	)
}

func startTracer(token string) error {
	lumigoToken := fmt.Sprintf("LumigoToken %s", token)
	ctx := context.Background()
	options := []otlptracehttp.Option{
		otlptracehttp.WithEndpoint(LUMIGO),
		otlptracehttp.WithHeaders(map[string]string{"Authorization": lumigoToken}),
	}

	client := otlptracehttp.NewClient(options...)
	exporter, err := otlptrace.New(ctx, client)
	if err != nil {
		return fmt.Errorf("error creating exporter %v", err)
	}

	traceProvider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(appResource()),
	)
	otel.SetTracerProvider(traceProvider)
	tracer = traceProvider.Tracer("tracer")

	return nil
}

func main() {
	fmt.Println("Starting")
	token, flag := os.LookupEnv("LUMIGO_TOKEN")
	if !flag || token == "" {
		log.Println("LUMIGO_TOKEN not set. Traces will not be received")
	}

	redisConnect()
	err := startTracer(token)
	if err != nil {
		log.Println("Faied to start tracer. You are flying blind. Good luck")
	}

	r := mux.NewRouter()
	r.Use(otelmux.Middleware("go-lumigo"))
	r.HandleFunc("/", handleHome)
	r.HandleFunc("/health", handleHealth)
	r.HandleFunc("/add/{key}/{value}", handleAdd)
	r.HandleFunc("/fetch/{key}", handleFetch)
	r.HandleFunc("/list", handleList)
	r.HandleFunc("/del/{key}", handleDelete)
	http.Handle("/", r)
	err = http.ListenAndServe(":8080", r)
	log.Fatal(err)
}
