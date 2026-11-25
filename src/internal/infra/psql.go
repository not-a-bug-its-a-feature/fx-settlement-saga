package infra

import (
    "context"
    "database/sql"
    "fmt"
    "time"

    _ "github.com/lib/pq"
)

// Psql Connection Poolを設定Pingで連携可能状態検証
func NewPostgres(ctx context.Context, dsn string) (*sql.DB, error) {
    db, err := sql.Open("postgres", dsn)

    if err != nil {
        return nil, fmt.Errorf("Open postgres: %w", err)
    }

    db.SetMaxOpenConns(20)
    db.SetMaxIdleConns(10)
    db.SetConnMaxLifetime(30 * time.Minute)

    pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)

    defer cancel()

    if err := db.PingContext(pingCtx); err != nil {
        _ = db.Close()

        return nil, fmt.Errorf("ping postgres: %w", err)
    }

    return db, nil
}