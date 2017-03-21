package main

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

var client = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "",
	DB:       1,
})

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	size, err := client.DbSize().Result()
	if err != nil {
		log.Fatal("Redis.Get: ", err)
	} else {
		fmt.Fprintf(w, `<html>
<head><title>URL Shorten service</title></head>
<body>We are serving %d URLs. Amazing!</body>
</html>
`, size)
	}
}

func RedirectHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	url, err := client.Get(key).Result()

	if err == redis.Nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 - Not Found"))
	} else if err != nil {
		log.Fatal("Redis.Get: ", err)
	} else {
		http.Redirect(w, r, url, 302)
	}
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", HomeHandler)
	//router.HandleFunc("/new", NewRedirectHandler)
	router.HandleFunc("/{key}", RedirectHandler)

	http.Handle("/", router)

	err := http.ListenAndServe("127.0.0.1:8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
