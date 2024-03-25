package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/alicebob/miniredis"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

func jsonToBidRequest(jsonData *string, request *BidRequest) {
	err := json.Unmarshal([]byte((*jsonData)), request)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
}

var ctx = context.Background()
var readWriteMutex sync.RWMutex

func startMiniRedis() *redis.Client {
	var rdb *redis.Client
	srv, err := miniredis.Run()
	if err != nil {
		fmt.Println("Error starting miniredis server:", err)
		return rdb
	}
	//defer srv.Close()
	rdb = redis.NewClient(&redis.Options{
		Addr: srv.Addr(),
	})
	return rdb
}
func pushToCache(json string, rdb *redis.Client, rwm *sync.RWMutex) {
	rwm.Lock()
	err := rdb.LPush("BidRequestQueue", json).Err()
	if err != nil {
		fmt.Println("Error pushing item to queue:", err)
		return
	}
	rwm.Unlock()
}
func popFromCache(rdb *redis.Client, rwm *sync.RWMutex) string {
	var response string
	rwm.Lock()
	if rdb.LLen("BidRequestQueue").Val() > 0 {
		val, err := rdb.RPop("BidRequestQueue").Result()
		if err != nil {
			fmt.Println("Error pushing item to queue:", err)

		}
		response = val
	}
	rwm.Unlock()
	return response
}

var miniRedis = startMiniRedis()

// go run DataSource.go main.go BidRequest.go AdCampaign.go
func main() {

	r := gin.Default()

	r.GET("/ui/live", func(c *gin.Context) {
		c.File("./bid.html")
	})
	r.GET("/bid/response", func(c *gin.Context) {
		var exchange string = adExchange()
		fmt.Println(exchange)
		c.Header("Content-Type", "application/json")
		c.JSON(http.StatusOK, exchange)
	})
	r.POST("/bid/request", func(c *gin.Context) {

		var bidRequest BidRequest
		if err := c.ShouldBindJSON(&bidRequest); err != nil {
			c.JSON(400, gin.H{"error": "Invalid JSON format"})
			return
		}
		json, err := json.Marshal(bidRequest)
		pushToCache(string(json), miniRedis, &readWriteMutex)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to marshal BidRequest"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "BidRequest received"})
	})

	if err := r.Run(":5758"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

/*
go func(msg string) {
	time.Sleep(time.Second * 3)
	for miniRedis.LLen("myQueue").Val() > 0 {
		val := popFromCache(miniRedis, &readWriteMutex)
		fmt.Println("cache length:", miniRedis.LLen("myQueue").Val(), "Popped item:", val)
		time.Sleep(time.Second * 2)
	}
	fmt.Println(msg)
}("going")
*/
