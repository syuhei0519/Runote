package redis

import (
    "context"
    "github.com/redis/go-redis/v9"
    "log"
    "os"
)

var Client *redis.Client
var Ctx = context.Background()

func InitRedis() {
	url := os.Getenv("REDIS_URL")
	if url == "" {
		log.Fatal("❌ REDIS_URL が設定されていません")
	}
	log.Println("url =", url)

	opt, err := redis.ParseURL(url)
	if err != nil {
		log.Fatalf("❌ REDIS_URL のパースに失敗しました: %v", err)
	}
	log.Printf("🔧 パース後の redis.Options: Addr=%s, DB=%d\n", opt.Addr, opt.DB)

	Client = redis.NewClient(opt)

	if _, err := Client.Ping(Ctx).Result(); err != nil {
		log.Fatalf("❌ Redis に接続できません: %v", err)
	}

	log.Println("✅ Redis に接続しました")
}

func FlushAll() {
    Client.FlushDB(context.Background())
}