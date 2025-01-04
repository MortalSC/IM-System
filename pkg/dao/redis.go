package dao

import (
	"context"
	"github.com/MortalSC/IM-System/config"
	"github.com/go-redis/redis/v8"
	"time"
)

// Rc : 全局变量 Rc，表示 RedisCache 的实例
var Rc *RedisCache

// RedisCache : 封装 redis.Client，用于操作 Redis 缓存
type RedisCache struct {
	rdb *redis.Client
}

// init 函数用于初始化 Redis 连接
func init() {
	// 初始化连接
	rdb := redis.NewClient(config.C.InitRedisOptions())
	// 连接健康检查
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		panic("Failed to connect to Redis: " + err.Error())
	}
	Rc = &RedisCache{rdb: rdb}
}

// Put 方法用于将 key-value 数据存入 Redis，并设置过期时间
// 参数：
// - ctx: 上下文，用于控制请求的生命周期
// - key: Redis 中的键
// - val: Redis 中的值
// - expire: 数据的过期时间
// 返回值：
// - error: 如果操作失败，返回错误信息；成功则返回 nil
func (rc *RedisCache) Put(ctx context.Context, key, val string, expire time.Duration) error {
	err := rc.rdb.Set(ctx, key, val, expire)
	if err != nil {
		return err.Err()
	}
	return nil
}

// Get 方法用于从 Redis 中获取指定 key 的值
// 参数：
// - ctx: 上下文，用于控制请求的生命周期
// - key: 要获取的键
// 返回值：
// - string: 获取到的值
// - error: 如果操作失败，返回错误信息；成功则返回 nil
func (rc *RedisCache) Get(ctx context.Context, key string) (string, error) {
	res, err := rc.rdb.Get(ctx, key).Result()

	return res, err
}
