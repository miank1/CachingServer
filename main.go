package main

import (
	"cache-server/cache"
	"encoding/json"
	"net/http"
	"strconv"
)

var c = cache.NewCache()

func main() {

	http.HandleFunc("/set", setHandler)
	http.HandleFunc("/get", getHandler)
	http.HandleFunc("/delete", deleteHandler)

	http.ListenAndServe(":8080", nil)
}

func setHandler(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	value := r.URL.Query().Get("value")
	ttl := r.URL.Query().Get("ttl")

	ttlSeconds := 60 // default

	if ttl != "" {
		parsed, err := strconv.Atoi(ttl)
		if err == nil {
			ttlSeconds = parsed
		}
	}

	c.Set(key, value, ttlSeconds)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("OK"))
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")

	val, exists := c.Get(key)
	if !exists {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"value": val})
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	c.Delete(key)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Deleted"))
}
