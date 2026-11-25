package infra

import (
    "context"
    "fmt"
    "time"

    "github.com/redis/go-redis/v9"
)

// Redis Clientを生成Pingでサーバ応答を確認
func NewRedis(ctx context.Context, addr string) (*redis.Client, error) {
    if addr == "" {
        return nil, fmt.Errorf("redis addr is required")
    }

    client := redis.NewClient(&redis.Options{
        Addr:   addr,
        DialTimeout: 3 * time.Second,
        ReadTimeout: 2 * time.Second,
        WriteTimeout: 2 * time.Second,
    })

    pingCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
    defer cancel()

    if err := client.Ping(pingCtx).Err(); err != nil {
        _ = client.Close()
        return nil, fmt.Errorf("ping redis: %w", err)
    }

    return client, nil
}