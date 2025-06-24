package main

import (
	"context"
	"log"
	"net/http"

	"github.com/devdevaraj/bender/handler"
	"github.com/devdevaraj/bender/init_redis" 
	"github.com/gorilla/mux"
	"github.com/redis/go-redis/v9" 
) 
        
var (  
	ctx = context.Background()
	rdb *redis.Client
)   
     
func main() {
	rdb = init_redis.InitRedis(ctx)
	router := mux.NewRouter()

	router.HandleFunc("/bridges/{image}", func(w http.ResponseWriter, r *http.Request) {
		handler.CreateBridge(w, r, rdb, ctx)
	}).Methods("POST")
	router.HandleFunc("/bridges/{name}", func(w http.ResponseWriter, r *http.Request) {
		handler.DeleteBridge(w, r, rdb, ctx)
	}).Methods("DELETE")
	router.HandleFunc("/bridges", handler.ListBridges).Methods("GET")

	log.Println("Starting Bender API server on :8080...")
	log.Fatal(http.ListenAndServe(":8080", router))
}
