package models

import (
	"time"

	"github.com/pgvector/pgvector-go"
)

type DocumentChunk struct {
	ID         string          `json:"id"`
	DocumentID string          `json:"document_id"`
	ChunkIndex int             `json:"chunk_index"`
	Content    string          `json:"content"`
	Embedding  pgvector.Vector `json:"embedding,omitempty"`
	CreatedAt  time.Time       `json:"created_at"`
}
