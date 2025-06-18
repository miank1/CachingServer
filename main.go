package main

import (
	"cache-server/cache"
	"encoding/json"
	"net/http"
)

var c = cache.NewCache()

func main() {

	http.HandleFunc("/set", setHandler)
	http.HandleFunc("/get", getHandler)
	http.HandleFunc("/delete", deleteHandler)
	http.HandleFunc("/stats", statsHandler)

	http.ListenAndServe(":8080", nil)
}

type setRequest struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	TTL   int    `json:"ttl"` // in seconds
}

type deleteRequest struct {
	Key string `json:"key"`
}

func setHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	var req setRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.Key == "" || req.Value == "" {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if req.TTL <= 0 {
		req.TTL = 60 // default TTL
	}

	c.Set(req.Key, req.Value, req.TTL)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("OK"))
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")

	val, exists := c.Get(key)
	if !exists {
		http.Error(w, "Key not found or expired", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"value": val})
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodDelete {
		http.Error(w, "Invalid Method", http.StatusMethodNotAllowed)
		return
	}

	var req deleteRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.Key == "" {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	c.Delete(req.Key)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Deleted"))
}

func statsHandler(w http.ResponseWriter, r *http.Request) {
	stats := c.Stats()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}
