package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

var client *redis.Client

func main() {
	InitRedis()
	api := gin.New()
	api.GET("/take", TakeCoupon)
	api.Run("127.0.0.1:8080")
}

func TakeCoupon(ctx *gin.Context) {
	const couponKey = "coupon"
	const lockKey = "coupon_lock"

	timeout := 30000
	waitTime := 50
	for {
		if timeout < 1 {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "timeout"})
			return
		}
		ok := client.SetNX(ctx, lockKey, "locked", time.Second).Val()
		if ok {
			defer client.Del(ctx, lockKey)
			break
		}
		timeout -= waitTime
		time.Sleep(time.Millisecond * time.Duration(waitTime))
		continue
	}

	n, err := client.Get(ctx, couponKey).Int()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	if n < 1 {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "库存不足"})
		return
	}
	err = client.Decr(ctx, couponKey).Err()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": "ok"})
}

func InitRedis() {
	client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
}
