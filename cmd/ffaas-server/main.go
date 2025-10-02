package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	// comentar esta si no vas a usar godotenv
	// "github.com/joho/godotenv"

	"github.com/redis/go-redis/v9"

	"github.com/Franconl/ffaas/internal/repo"
	"github.com/Franconl/ffaas/internal/repo/cached"
	"github.com/Franconl/ffaas/internal/repo/memory"
	"github.com/Franconl/ffaas/internal/repo/postgres"

	_ "github.com/jackc/pgx/v5/stdlib" // driver para sql.Open("pgx", ...)
)

func main() {
	// _ = godotenv.Load() // opcional si usás .env

	useMemory := os.Getenv("USE_MEMORY") == "true"

	var store repo.Flags

	if useMemory {
		// 🔹 Repositorio en memoria (ideal para dev rápido)
		store = memory.New()
		log.Println("⚡ Usando repositorio en memoria")
	} else {
		// 🔹 Config DB
		pgUser := getEnv("DB_USER", "app")
		pgPass := getEnv("DB_PASS", "app")
		pgHost := getEnv("DB_HOST", "localhost")
		pgPort := getEnv("DB_PORT", "5432")
		pgName := getEnv("DB_NAME", "appdb")

		dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
			pgUser, pgPass, pgHost, pgPort, pgName,
		)

		// Conexión a Postgres
		db, err := sql.Open("pgx", dsn)
		if err != nil {
			log.Fatal("❌ Error conectando a Postgres:", err)
		}
		defer db.Close()

		if err := db.Ping(); err != nil {
			log.Fatal("❌ Error al hacer ping a Postgres:", err)
		}
		pgRepo := postgres.New(db)

		// 🔹 Redis (opcional)
		redisAddr := getEnv("REDIS_ADDR", "localhost:6379")
		rdb := redis.NewClient(&redis.Options{
			Addr:     redisAddr,
			Password: os.Getenv("REDIS_PASS"),
			DB:       0,
		})
		if err := rdb.Ping(context.Background()).Err(); err != nil {
			log.Fatal("❌ Error conectando a Redis:", err)
		}

		// Repo cacheado (Postgres + Redis)
		store = cached.New(pgRepo, rdb, 60*time.Second)
		log.Println("⚡ Usando Postgres + Redis")
	}

	// --- HTTP Router ---
	r := chi.NewRouter()

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	r.Get("/flags", func(w http.ResponseWriter, r *http.Request) {
		flags, err := store.List()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		writeJSON(w, flags)
	})

	// --- Server ---
	addr := ":" + getEnv("APP_PORT", "8080")
	log.Println("🚀 API escuchando en", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}

// Helpers -----------------

func getEnv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	_ = enc.Encode(v)
}
