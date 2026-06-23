package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"browser-copilot-backend/internal/config"
	"browser-copilot-backend/internal/db"
)

func main() {
	// Initialize configurations and DBs
	config.LoadConfig()
	db.InitPostgres()
	db.InitRedis()

	ctx := context.Background()

	// 1. Setup dummy data
	key := "benchmark_session_123"
	payload := `[{"role":"user","content":"Hello AI"},{"role":"ai","content":"How can I help?"}]`

	// Insert into Redis
	err := db.RDB.Set(ctx, key, payload, 0).Err()
	if err != nil {
		log.Fatal("Failed to set Redis:", err)
	}

	// Insert into Postgres (We'll use a temporary table or just benchmark a simple query if we don't have a table ready)
	_, err = db.PG.Exec("CREATE TABLE IF NOT EXISTS benchmark_data (id VARCHAR PRIMARY KEY, data TEXT)")
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.PG.Exec("INSERT INTO benchmark_data (id, data) VALUES ($1, $2) ON CONFLICT (id) DO UPDATE SET data = $2", key, payload)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("--- Running Benchmarks (10 iterations each) ---")

	// 2. Benchmark Postgres
	var totalPgTime time.Duration
	for i := 0; i < 10; i++ {
		start := time.Now()
		var res string
		err := db.PG.QueryRow("SELECT data FROM benchmark_data WHERE id = $1", key).Scan(&res)
		if err != nil {
			log.Fatal(err)
		}
		totalPgTime += time.Since(start)
	}
	avgPg := totalPgTime / 10

	// 3. Benchmark Redis
	var totalRedisTime time.Duration
	for i := 0; i < 10; i++ {
		start := time.Now()
		_, err := db.RDB.Get(ctx, key).Result()
		if err != nil {
			log.Fatal(err)
		}
		totalRedisTime += time.Since(start)
	}
	avgRedis := totalRedisTime / 10

	fmt.Printf("PostgreSQL Average Latency: %v\n", avgPg)
	fmt.Printf("Redis Average Latency: %v\n", avgRedis)

	// 4. Do the math
	diff := float64(avgPg.Nanoseconds()-avgRedis.Nanoseconds()) / float64(avgPg.Nanoseconds()) * 100
	fmt.Printf("\nRedis reduced latency by: %.2f%%\n", diff)

	// Cleanup
	db.PG.Exec("DROP TABLE benchmark_data")
	db.RDB.Del(ctx, key)
}
