package db

import (
	"context"
	"database/sql"
	"log/slog"
	"os"

	"browser-copilot-backend/internal/config"
	_ "github.com/lib/pq"
	"github.com/pgvector/pgvector-go"
)

var PG *sql.DB

func InitPostgres() {
	dbURL := config.AppConfig.DBURL

	var err error
	PG, err = sql.Open("postgres", dbURL)
	if err != nil {
		slog.Error("Failed to open Postgres connection", "error", err)
		os.Exit(1)
	}

	if err = PG.Ping(); err != nil {
		slog.Error("Failed to ping Postgres", "error", err)
		os.Exit(1)
	}

	slog.Info("🐘 PostgreSQL Connected successfully!")
}

// StoreMessage stores a chat message in the DB
func StoreMessage(ctx context.Context, sessionID, role, content string) error {
	// First ensure conversation exists
	_, err := PG.ExecContext(ctx, "INSERT INTO conversations (id) VALUES ($1) ON CONFLICT (id) DO NOTHING", sessionID)
	if err != nil {
		return err
	}

	_, err = PG.ExecContext(ctx, "INSERT INTO messages (conversation_id, role, content) VALUES ($1, $2, $3)", sessionID, role, content)
	return err
}

// SaveDocumentAndEmbeddings saves a webpage and its vector embeddings
func SaveDocumentAndEmbeddings(ctx context.Context, url, title, content string, chunks []string, embeddings [][]float32) error {
	tx, err := PG.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var docID int
	err = tx.QueryRowContext(ctx, `
		INSERT INTO webpage_documents (url, title, content) 
		VALUES ($1, $2, $3) 
		ON CONFLICT (url) DO UPDATE SET title = EXCLUDED.title, content = EXCLUDED.content
		RETURNING id
	`, url, title, content).Scan(&docID)
	if err != nil {
		return err
	}

	// Delete old embeddings for this doc
	_, err = tx.ExecContext(ctx, "DELETE FROM embeddings WHERE document_id = $1", docID)
	if err != nil {
		return err
	}

	for i, chunk := range chunks {
		vec := pgvector.NewVector(embeddings[i])
		_, err = tx.ExecContext(ctx, "INSERT INTO embeddings (document_id, chunk_text, embedding) VALUES ($1, $2, $3)", docID, chunk, vec)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// SearchSimilarChunks finds similar texts using pgvector
func SearchSimilarChunks(ctx context.Context, queryEmbedding []float32, limit int) ([]string, error) {
	vec := pgvector.NewVector(queryEmbedding)
	rows, err := PG.QueryContext(ctx, `
		SELECT chunk_text 
		FROM embeddings 
		ORDER BY embedding <-> $1 
		LIMIT $2
	`, vec, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chunks []string
	for rows.Next() {
		var chunk string
		if err := rows.Scan(&chunk); err != nil {
			return nil, err
		}
		chunks = append(chunks, chunk)
	}
	return chunks, nil
}
