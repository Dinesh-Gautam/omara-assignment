package services

import (
	"database/sql"
	"fmt"
	"log"
	"mime/multipart"
	"strategic-insight-analyst/backend/config"
	"strategic-insight-analyst/backend/database"
	"strategic-insight-analyst/backend/internal/llm"
	"strategic-insight-analyst/backend/internal/processor"
	"strategic-insight-analyst/backend/internal/storage"
	"strategic-insight-analyst/backend/models"
	"strings"
	"time"

	"github.com/lib/pq"
	"github.com/pgvector/pgvector-go"
	uuid "github.com/satori/go.uuid"
)

func ProcessAndSaveDocument(file multipart.File, handler *multipart.FileHeader, userID string) (models.Document, error) {
	gcsPath, err := storage.UploadFile(file, handler)
	if err != nil {
		return models.Document{}, fmt.Errorf("failed to upload file to GCS: %w", err)
	}

	doc := models.Document{
		ID:        uuid.NewV4().String(),
		UserID:    userID,
		FileName:  handler.Filename,
		GCSPath:   gcsPath,
		Status:    "processing",
		CreatedAt: time.Now(),
	}
	query := `INSERT INTO documents (id, user_id, file_name, gcs_path, status, created_at) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err = database.DB.Exec(query, doc.ID, doc.UserID, doc.FileName, doc.GCSPath, doc.Status, doc.CreatedAt)
	if err != nil {
		return models.Document{}, fmt.Errorf("failed to create document record: %w", err)
	}

	go func() {
		file.Seek(0, 0)
		fileType := handler.Header.Get("Content-Type")
		var processingErr error

		if fileType == "application/pdf" {
			log.Printf("Starting PDF chunk processing for document %s", doc.ID)
			processingErr = processor.ProcessPDFChunks(file, 10000, 200, func(chunk string, chunkIndex int) error {
				return processChunk(chunk, doc.ID, chunkIndex)
			})
		} else {
			textContent, err := processor.ExtractText(file, fileType)
			if err != nil {
				processingErr = fmt.Errorf("failed to extract text: %w", err)
			} else {
				chunks := processor.ChunkText(textContent, 10000, 500)

				for i, chunk := range chunks {
					if err := processChunk(chunk, doc.ID, i); err != nil {
						processingErr = fmt.Errorf("failed to process chunk %d: %w", i, err)
						break
					}
				}
			}
		}

		var finalStatus string
		var finalError string
		if processingErr != nil {
			log.Printf("ERROR: Failed to process document %s: %v", doc.ID, processingErr)
			finalStatus = "failed"
			finalError = processingErr.Error()
		} else {
			log.Printf("Successfully processed all chunks for document %s", doc.ID)
			finalStatus = "processed"
			finalError = ""
		}

		updateQuery := `UPDATE documents SET status = $1, processing_error = $2 WHERE id = $3`
		_, err := database.DB.Exec(updateQuery, finalStatus, finalError, doc.ID)
		if err != nil {
			log.Printf("ERROR: Failed to update document status for %s: %v", doc.ID, err)
		}
	}()

	return doc, nil
}

func processChunk(chunk string, docID string, chunkIndex int) error {
	log.Printf("Processing chunk %d for document %s", chunkIndex, docID)

	embedding, err := llm.GetEmbedding(chunk)
	if err != nil {
		log.Printf("ERROR: Failed to generate embedding for chunk %d for document %s: %v", chunkIndex, docID, err)
		return err
	}

	chunkModel := models.DocumentChunk{
		ID:         uuid.NewV4().String(),
		DocumentID: docID,
		ChunkIndex: chunkIndex,
		Content:    chunk,
		Embedding:  pgvector.NewVector(embedding),
		CreatedAt:  time.Now(),
	}
	query := `INSERT INTO document_chunks (id, document_id, chunk_index, content, embedding, created_at) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err = database.DB.Exec(query, chunkModel.ID, chunkModel.DocumentID, chunkModel.ChunkIndex, chunkModel.Content, chunkModel.Embedding, chunkModel.CreatedAt)
	if err != nil {
		log.Printf("ERROR: Failed to save chunk %d for document %s: %v", chunkIndex, docID, err)
		return err
	}
	log.Printf("Successfully saved chunk %d for document %s", chunkIndex, docID)
	return nil
}

func GetUserDocuments(userID string) ([]models.Document, error) {
	query := `SELECT id, file_name, gcs_path, status, created_at FROM documents WHERE user_id = $1 ORDER BY created_at DESC`
	rows, err := database.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	documents := make([]models.Document, 0)
	for rows.Next() {
		var doc models.Document
		if err := rows.Scan(&doc.ID, &doc.FileName, &doc.GCSPath, &doc.Status, &doc.CreatedAt); err != nil {
			return nil, err
		}
		doc.UserID = userID
		documents = append(documents, doc)
	}

	return documents, nil
}

func GetDocumentStatus(documentID, userID string) (models.Document, error) {
	var doc models.Document
	var processingError sql.NullString
	query := `SELECT id, file_name, gcs_path, status, created_at, processing_error FROM documents WHERE id = $1 AND user_id = $2`
	err := database.DB.QueryRow(query, documentID, userID).Scan(&doc.ID, &doc.FileName, &doc.GCSPath, &doc.Status, &doc.CreatedAt, &processingError)
	if err != nil {
		return models.Document{}, err
	}
	if processingError.Valid {
		doc.ProcessingError = processingError.String
	}
	doc.UserID = userID
	return doc, nil
}

func getObjectName(gcsPath string) string {
	if strings.HasPrefix(gcsPath, "gs://") {
		bucketName := config.AppConfig.GCSBucketName
		prefix := fmt.Sprintf("gs://%s/", bucketName)
		return strings.TrimPrefix(gcsPath, prefix)
	}
	return gcsPath
}

func DownloadDocument(documentID, userID string) ([]byte, string, error) {
	var gcsPath, fileName string
	query := `SELECT gcs_path, file_name FROM documents WHERE id = $1 AND user_id = $2`
	err := database.DB.QueryRow(query, documentID, userID).Scan(&gcsPath, &fileName)
	if err != nil {
		return nil, "", err
	}

	objectName := getObjectName(gcsPath)
	fileContent, err := storage.DownloadFile(objectName)
	if err != nil {
		return nil, "", err
	}

	return fileContent, fileName, nil
}

func DeleteDocument(documentID, userID string) error {
	var gcsPath string
	query := `SELECT gcs_path FROM documents WHERE id = $1 AND user_id = $2`
	err := database.DB.QueryRow(query, documentID, userID).Scan(&gcsPath)
	if err != nil {
		return err
	}

	objectName := getObjectName(gcsPath)
	err = storage.DeleteFile(objectName)
	if err != nil {
		return err
	}

	query = `DELETE FROM documents WHERE id = $1`
	_, err = database.DB.Exec(query, documentID)
	return err
}

func GetDocumentContent(documentID string) (string, error) {
	query := `
		SELECT content
		FROM document_chunks
		WHERE document_id = $1
		ORDER BY chunk_index
	`
	rows, err := database.DB.Query(query, documentID)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	var chunks []string
	for rows.Next() {
		var content string
		if err := rows.Scan(&content); err != nil {
			return "", err
		}
		chunks = append(chunks, content)
	}

	return strings.Join(chunks, "\n\n"), nil
}

func GetDocumentsByIDs(documentIDs []string, userID string) ([]models.Document, error) {
	if len(documentIDs) == 0 {
		return []models.Document{}, nil
	}

	query := `SELECT id, file_name, gcs_path, status, created_at FROM documents WHERE id = ANY($1) AND user_id = $2`
	rows, err := database.DB.Query(query, pq.Array(documentIDs), userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	documents := make([]models.Document, 0)
	for rows.Next() {
		var doc models.Document
		if err := rows.Scan(&doc.ID, &doc.FileName, &doc.GCSPath, &doc.Status, &doc.CreatedAt); err != nil {
			return nil, err
		}
		doc.UserID = userID
		documents = append(documents, doc)
	}

	return documents, nil
}
