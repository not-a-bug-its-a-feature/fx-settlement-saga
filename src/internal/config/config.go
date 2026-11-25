package config

import (
    "fmt"
    "os"
)

type Config struct {
    DBDSN           string
    RedisAddr       string
    KafkaBrokers    string
    HTTPAddr        string
}
// 環境変数を読み込みConfigを組み立てる
// 必須の値がなければエラー
func Load() (Config, error) {
    cfg := Config {
        DBDSN:  getEnv("DB_DSN", ""),
        RedisAddr:  getEnv("REDIS_ADDR", "localhost:6379"),
        KafkaBrokers: getEnv("KAFKA_BROKERS", "localhost:9092"),
        HTTPAddr: getEnv("HTTP_ADDR", ":8080"),
    }

    if cfg.DBDSN == "" {
        return Config{}, fmt.Errorf("DB_DSN is required")
    }

    return cfg, nil
}

// 指定の環境変数を探し値がなければ基本値（fallback）
func getEnv(key, fallback string) string {
    if v := os.Getenv(key); v != "" {
        return v
    }

    return fallback
}