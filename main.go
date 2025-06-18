package main

import (
	"cache-server/handler"
	"net/http"
)

func main() {

	http.HandleFunc("/set", handler.SetHandler)
	http.HandleFunc("/get", handler.GetHandler)
	http.HandleFunc("/delete", handler.DeleteHandler)
	http.HandleFunc("/stats", handler.StatsHandler)

	http.ListenAndServe(":8080", nil)
}
