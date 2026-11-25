package main

import (
    "context"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/not-a-bug-its-a-feature/fx-settlement-saga/internal/config"
    "github.com/not-a-bug-its-a-feature/fx-settlement-saga/internal/infra"
)

func main() {
    ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
    defer stop()

    cfg, err := config.Load()
    if err != nil {
        log.Fatal("load config: %v", err)
    }

    db, err := infra.NewPostgres(ctx, cfg.DBDSN)

    if err != nil {
        log.Fatal(err)
    }

    defer db.Close()

    cache, err := infra.NewRedis(ctx, cfg.RedisAddr)
    if err != nil {
        log.Fatal("init redis %v", err)
    }
    defer cache.Close()

    kafkaWriter, err := infra.NewKafkaWriter(cfg.KafkaBrokers, "saga.events")
    if err != nil {
        log.Fatal("init kafka writer: %v", err)
    }
    defer kafkaWriter.Close()

    mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	srv := &http.Server{
		Addr:    cfg.HTTPAddr,
		Handler: mux,
	}

    go func() {
		log.Printf("api server listening on %s", cfg.HTTPAddr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen and serve: %v", err)
		}
	}()

    <-ctx.Done()
	log.Println("shutdown signal received, closing server...")

    shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("server shutdown error: %v", err)
	}
}