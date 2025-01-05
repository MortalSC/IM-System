package redis

import (
	"context"
	"github.com/MortalSC/IM-System/lib/cache"
	"time"
)
import "github.com/go-redis/redis/v8"

type IMRedisCache struct {
	rdb *redis.Client
}

// NewRedisCache 创建一个新的 Redis 缓存实例
// 参数：
// - addr: Redis 地址
// - password: Redis 密码
// - db: Redis 数据库索引
// 返回值：
// - cache.Cache: 实现了通用 Cache 接口的 Redis 缓存实例
// - error: 如果初始化失败，返回错误
func NewRedisCache(addr, password string, db int) (cache.Cache, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	// 检查连接是否正常
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	return &IMRedisCache{rdb: rdb}, nil
}

// Put 方法用于将 key-value 数据存入 Redis，并设置过期时间
// 参数：
// - ctx: 上下文，用于控制请求的生命周期
// - key: Redis 中的键
// - val: Redis 中的值
// - expire: 数据的过期时间
// 返回值：
// - error: 如果操作失败，返回错误信息；成功则返回 nil
func (rc *IMRedisCache) Put(ctx context.Context, key, val string, expire time.Duration) error {
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
func (rc *IMRedisCache) Get(ctx context.Context, key string) (string, error) {
	res, err := rc.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		// 返回空字符串表示 key 不存在
		return "", nil
	}
	return res, err
}
