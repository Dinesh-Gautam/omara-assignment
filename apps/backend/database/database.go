package database

import (
	"database/sql"
	"fmt"
	"log"
	"strategic-insight-analyst/backend/config"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Connect() {
	var err error
	dsn := config.AppConfig.GetDBDSN()

	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	fmt.Println("Successfully connected to database")
}

func Migrate() {
	if _, err := DB.Exec("CREATE EXTENSION IF NOT EXISTS vector;"); err != nil {
		log.Fatal("Failed to create vector extension:", err)
	}

	usersQuery := `
	   CREATE TABLE IF NOT EXISTS users (
	       id VARCHAR(255) PRIMARY KEY, -- Firebase Auth UID
	       email VARCHAR(255) NOT NULL UNIQUE,
	       auth_method VARCHAR(50),
	       created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	   );`

	if _, err := DB.Exec(usersQuery); err != nil {
		log.Fatal("Failed to create users table:", err)
	}

	documentsQuery := `
	   CREATE TABLE IF NOT EXISTS documents (
	       id VARCHAR(255) PRIMARY KEY,
	       user_id VARCHAR(255) NOT NULL,
	       file_name VARCHAR(255) NOT NULL,
	       gcs_path VARCHAR(255) NOT NULL,
	       status VARCHAR(50) NOT NULL DEFAULT 'processing',
		   processing_error TEXT,
	       created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
	       FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	   );`

	if _, err := DB.Exec(documentsQuery); err != nil {
		log.Fatal("Failed to create documents table:", err)
	}

	documentChunksQuery := `
	   CREATE TABLE IF NOT EXISTS document_chunks (
	       id VARCHAR(255) PRIMARY KEY,
	       document_id VARCHAR(255) NOT NULL,
	    chunk_index INTEGER NOT NULL,
	       content TEXT NOT NULL,
	       embedding vector(768),
	       created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
	       FOREIGN KEY (document_id) REFERENCES documents(id) ON DELETE CASCADE
	   );`

	if _, err := DB.Exec(documentChunksQuery); err != nil {
		log.Fatal("Failed to create document_chunks table:", err)
	}

	chatHistoryQuery := `
	   CREATE TABLE IF NOT EXISTS chat_history (
	       id VARCHAR(255) PRIMARY KEY, -- Unique ID for the chat message
	   document_id VARCHAR(255) NOT NULL,
	   user_id VARCHAR(255) NOT NULL,
	   message_type VARCHAR(255) NOT NULL, -- 'user' for query, 'ai' for response
	   message_content TEXT NOT NULL,
	   attached_documents JSONB,
	   timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	   FOREIGN KEY (document_id) REFERENCES documents(id) ON DELETE CASCADE,
	   FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	   );`

	_, err := DB.Exec(chatHistoryQuery)
	if err != nil {
		log.Fatal("Failed to create chat_history table:", err)
	}

	if err := createIndexes(); err != nil {
		log.Fatal("Failed to create indexes: ", err)
	}

	fmt.Println("Database migration completed")
}

func createIndexes() error {
	indexQueries := []string{
		"CREATE INDEX IF NOT EXISTS idx_documents_user_id ON documents (user_id);",
		"CREATE INDEX IF NOT EXISTS idx_document_chunks_document_id ON document_chunks (document_id);",
		"CREATE INDEX IF NOT EXISTS idx_documents_status ON documents (status);",
		"CREATE INDEX IF NOT EXISTS idx_document_chunks_embedding ON document_chunks USING hnsw (embedding vector_l2_ops);",
		"CREATE INDEX IF NOT EXISTS idx_chat_history_document_user ON chat_history (document_id, user_id);",
	}

	for _, query := range indexQueries {
		_, err := DB.Exec(query)
		if err != nil {
			return fmt.Errorf("failed to create index with query '%s': %w", query, err)
		}
	}

	log.Println("Database indexes created successfully.")
	return nil
}
