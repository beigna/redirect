package main

import (
    //"fmt"
    "log"
    "net/http"
    "github.com/go-redis/redis"
)

type redirect struct {
    key string
    url string
}

func handler(w http.ResponseWriter, r *http.Request) {
    client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       1,  // use default DB
	})

    key := r.URL.Path[1:]
    val, err := client.Get(key).Result()

    if err == redis.Nil {
        w.WriteHeader(http.StatusNotFound)
        w.Write([]byte("404 - Not Found"))
    } else if err != nil {
        log.Fatal("Redis.Get: ", err)
    } else {
        http.Redirect(w, r, val, 302)
    }
}

func main() {
    http.HandleFunc("/", handler)
    err := http.ListenAndServe(":8080", nil)
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}
