package database

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"tabimoney/internal/config"

	"github.com/redis/go-redis/v9"
)

var (
	RedisClient *redis.Client

	rateLimitScript = redis.NewScript(`
local current = redis.call("INCR", KEYS[1])
if current == 1 then
  redis.call("EXPIRE", KEYS[1], ARGV[1])
end
return current
`)
)

func InitRedis(cfg *config.Config) error {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.GetRedisAddr(),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return fmt.Errorf("failed to connect to Redis: %w", err)
	}

	RedisClient = rdb
	log.Println("Redis connection established successfully")
	return nil
}

func CloseRedis() error {
	if RedisClient == nil {
		return nil
	}

	return RedisClient.Close()
}

// GetRedisClient returns the Redis client instance
func GetRedisClient() *redis.Client {
	return RedisClient
}

// Cache operations
func SetCache(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return RedisClient.Set(ctx, key, value, expiration).Err()
}

func GetCache(ctx context.Context, key string) (string, error) {
	return RedisClient.Get(ctx, key).Result()
}

func DeleteCache(ctx context.Context, key string) error {
	return RedisClient.Del(ctx, key).Err()
}

func ExistsCache(ctx context.Context, key string) (bool, error) {
	result, err := RedisClient.Exists(ctx, key).Result()
	return result > 0, err
}

// Session management
func SetSession(ctx context.Context, sessionID string, userID uint64, expiration time.Duration) error {
	key := fmt.Sprintf("session:%s", sessionID)
	return RedisClient.Set(ctx, key, userID, expiration).Err()
}

func GetSession(ctx context.Context, sessionID string) (uint64, error) {
	key := fmt.Sprintf("session:%s", sessionID)
	result, err := RedisClient.Get(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	
	var userID uint64
	_, err = fmt.Sscanf(result, "%d", &userID)
	return userID, err
}

func DeleteSession(ctx context.Context, sessionID string) error {
	key := fmt.Sprintf("session:%s", sessionID)
	return RedisClient.Del(ctx, key).Err()
}

// Rate limiting
func IncrementRateLimit(ctx context.Context, key string, window time.Duration) (int64, error) {
	if RedisClient == nil {
		return 0, fmt.Errorf("Redis not initialized")
	}

	windowSeconds := int64(window / time.Second)
	if windowSeconds <= 0 {
		windowSeconds = 1
	}

	result, err := rateLimitScript.Run(ctx, RedisClient, []string{key}, windowSeconds).Result()
	if err != nil {
		return 0, err
	}

	switch v := result.(type) {
	case int64:
		return v, nil
	case string:
		return strconv.ParseInt(v, 10, 64)
	default:
		return 0, fmt.Errorf("unexpected rate limit script result: %T", v)
	}
}

func GetRateLimit(ctx context.Context, key string) (int64, error) {
	return RedisClient.Get(ctx, key).Int64()
}

// Dashboard cache
func SetDashboardCache(ctx context.Context, userID uint64, period string, data interface{}, expiration time.Duration) error {
	key := fmt.Sprintf("dashboard:%d:%s", userID, period)
	return SetCache(ctx, key, data, expiration)
}

func GetDashboardCache(ctx context.Context, userID uint64, period string) (string, error) {
	key := fmt.Sprintf("dashboard:%d:%s", userID, period)
	return GetCache(ctx, key)
}

func DeleteDashboardCache(ctx context.Context, userID uint64) error {
	pattern := fmt.Sprintf("dashboard:%d:*", userID)
	return deleteByPattern(ctx, pattern)
}

// AI analysis cache
func SetAIAnalysisCache(ctx context.Context, userID uint64, analysisType string, data interface{}, expiration time.Duration) error {
	key := fmt.Sprintf("ai_analysis:%d:%s", userID, analysisType)
	return SetCache(ctx, key, data, expiration)
}

func GetAIAnalysisCache(ctx context.Context, userID uint64, analysisType string) (string, error) {
	key := fmt.Sprintf("ai_analysis:%d:%s", userID, analysisType)
	return GetCache(ctx, key)
}

func DeleteAIAnalysisCache(ctx context.Context, userID uint64) error {
	pattern := fmt.Sprintf("ai_analysis:%d:*", userID)
	return deleteByPattern(ctx, pattern)
}

// Notification queue
func PushNotification(ctx context.Context, userID uint64, notification interface{}) error {
	key := fmt.Sprintf("notifications:%d", userID)
	return RedisClient.LPush(ctx, key, notification).Err()
}

func PopNotification(ctx context.Context, userID uint64) (string, error) {
	key := fmt.Sprintf("notifications:%d", userID)
	return RedisClient.RPop(ctx, key).Result()
}

func GetNotificationQueueLength(ctx context.Context, userID uint64) (int64, error) {
	key := fmt.Sprintf("notifications:%d", userID)
	return RedisClient.LLen(ctx, key).Result()
}

// Health check for Redis connection
func RedisHealthCheck() error {
	if RedisClient == nil {
		return fmt.Errorf("Redis not initialized")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		return fmt.Errorf("Redis ping failed: %w", err)
	}

	return nil
}

func deleteByPattern(ctx context.Context, pattern string) error {
	if RedisClient == nil {
		return fmt.Errorf("Redis not initialized")
	}

	var (
		cursor uint64
		batch  = make([]string, 0, 100)
	)

	for {
		keys, nextCursor, err := RedisClient.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			return err
		}

		if len(keys) > 0 {
			batch = append(batch, keys...)
			if len(batch) >= 100 {
				if err := RedisClient.Del(ctx, batch...).Err(); err != nil {
					return err
				}
				batch = batch[:0]
			}
		}

		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}

	if len(batch) > 0 {
		if err := RedisClient.Del(ctx, batch...).Err(); err != nil {
			return err
		}
	}

	return nil
}
