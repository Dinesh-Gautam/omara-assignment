package services

import (
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"strategic-insight-analyst/backend/database"
	"strategic-insight-analyst/backend/internal/llm"
	"strategic-insight-analyst/backend/models"
	"strings"
	"time"

	"github.com/pgvector/pgvector-go"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/genai"
)

func SaveUserMessage(documentID, userID, userMessage string, attachedDocs []models.Document) (models.ChatMessage, error) {
	var attachedDocsJSON []byte
	var err error

	// Always create a valid JSON array, even if it's empty.
	serializableDocs := make([]map[string]string, len(attachedDocs))
	for i, doc := range attachedDocs {
		serializableDocs[i] = map[string]string{"id": doc.ID, "title": doc.FileName}
	}
	attachedDocsJSON, err = json.Marshal(serializableDocs)
	if err != nil {
		return models.ChatMessage{}, fmt.Errorf("failed to marshal attached documents: %w", err)
	}

	message := models.ChatMessage{
		ID:                uuid.NewV4().String(),
		DocumentID:        documentID,
		UserID:            userID,
		MessageType:       "user",
		MessageContent:    userMessage,
		Timestamp:         time.Now(),
		AttachedDocuments: string(attachedDocsJSON),
	}

	query := `INSERT INTO chat_history (id, document_id, user_id, message_type, message_content, timestamp, attached_documents) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err = database.DB.Exec(query, message.ID, message.DocumentID, message.UserID, message.MessageType, message.MessageContent, message.Timestamp, message.AttachedDocuments)
	return message, err
}

func GetRelevantContext(documentID, userMessage string, attachedDocIDs []string, userID string) (string, error) {
	var contextBuilder strings.Builder

	// 1. Fetch content from attached documents
	if len(attachedDocIDs) > 0 {
		docs, err := GetDocumentsByIDs(attachedDocIDs, userID)
		if err != nil {
			return "", fmt.Errorf("failed to get attached documents: %w", err)
		}

		for _, doc := range docs {
			content, err := GetDocumentContent(doc.ID)
			if err != nil {
				// Log the error but continue, so one failed doc doesn't stop the whole process
				log.Printf("Warning: failed to get content for attached document %s: %v", doc.ID, err)
				continue
			}
			contextBuilder.WriteString(fmt.Sprintf("<document>\n<title>%s</title>\n<content>\n%s\n</content>\n</document>\n", doc.FileName, content))
		}
	}

	// 2. Fetch relevant chunks from the main document
	userMessageEmbedding, err := llm.GetEmbedding(userMessage)
	if err != nil {
		return "", err
	}

	query := `
		SELECT content, chunk_index
		FROM document_chunks
		WHERE document_id = $1
		ORDER BY embedding <=> $2
		LIMIT 3
	`
	rows, err := database.DB.Query(query, documentID, pgvector.NewVector(userMessageEmbedding))
	if err != nil {
		return "", err
	}
	defer rows.Close()

	type chunkWithIndex struct {
		Content    string
		ChunkIndex int
	}

	var chunksWithIndices []chunkWithIndex
	for rows.Next() {
		var c chunkWithIndex
		if err := rows.Scan(&c.Content, &c.ChunkIndex); err != nil {
			return "", err
		}
		chunksWithIndices = append(chunksWithIndices, c)
	}

	sort.Slice(chunksWithIndices, func(i, j int) bool {
		return chunksWithIndices[i].ChunkIndex < chunksWithIndices[j].ChunkIndex
	})

	var chunks []string
	for _, c := range chunksWithIndices {
		chunks = append(chunks, c.Content)
	}

	// 3. Combine the contexts
	mainDocContext := strings.Join(chunks, "\n\n")
	if contextBuilder.Len() > 0 {
		// We have attached documents, so we wrap the main doc context as well
		mainDocInfo, err := GetDocumentStatus(documentID, userID)
		if err != nil {
			log.Printf("Warning: could not get main document info for %s: %v", documentID, err)
			// Fallback to just using the content
			contextBuilder.WriteString(fmt.Sprintf("<document>\n<title>Current Document</title>\n<content>\n%s\n</content>\n</document>\n", mainDocContext))
		} else {
			contextBuilder.WriteString(fmt.Sprintf("<document>\n<title>%s</title>\n<content>\n%s\n</content>\n</document>\n", mainDocInfo.FileName, mainDocContext))
		}
		return contextBuilder.String(), nil
	}

	// If no attached docs, return only the relevant chunks from the main document
	return mainDocContext, nil
}

func GetChatHistory(documentID string) ([]models.ChatMessage, error) {
	query := `SELECT id, document_id, user_id, message_type, message_content, timestamp, attached_documents FROM chat_history WHERE document_id = $1 ORDER BY timestamp`
	rows, err := database.DB.Query(query, documentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	history := make([]models.ChatMessage, 0)
	for rows.Next() {
		var msg models.ChatMessage
		var attachedDocsBytes []byte
		if err := rows.Scan(&msg.ID, &msg.DocumentID, &msg.UserID, &msg.MessageType, &msg.MessageContent, &msg.Timestamp, &attachedDocsBytes); err != nil {
			return nil, err
		}
		// The frontend expects a JSON string, so we just assign it.
		// The model has `omitempty`, so it will be null if empty.
		msg.AttachedDocuments = string(attachedDocsBytes)
		history = append(history, msg)
	}
	return history, nil
}

func GetChatHistoryForLLM(documentID string) ([]*genai.Content, error) {
	query := `SELECT message_type, message_content FROM chat_history WHERE document_id = $1 ORDER BY timestamp`
	rows, err := database.DB.Query(query, documentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []*genai.Content
	for rows.Next() {
		var messageType, messageContent string
		if err := rows.Scan(&messageType, &messageContent); err != nil {
			return nil, err
		}
		role := genai.RoleUser
		if messageType == "ai" {
			role = genai.RoleModel
		}
		history = append(history, genai.NewContentFromText(messageContent, genai.Role(role)))
	}
	return history, nil
}

func StreamChatResponse(userMessage, contextText string, history []*genai.Content, streamChan chan string) (string, error) {
	hasAttachedDocs := strings.Contains(contextText, "<document>")
	fullResponse, err := llm.CallGeminiStream(userMessage, contextText, history, hasAttachedDocs, streamChan)
	if err != nil {
		log.Printf("Error from CallGeminiStream: %v", err)
		return "", err // Propagate the error
	}

	// The handler will be responsible for saving the AI message now.
	return fullResponse, nil
}

func SaveAIMessage(documentID, userID, aiResponse string) (models.ChatMessage, error) {
	message := models.ChatMessage{
		ID:             uuid.NewV4().String(),
		DocumentID:     documentID,
		UserID:         userID,
		MessageType:    "ai",
		MessageContent: aiResponse,
		Timestamp:      time.Now(),
	}
	query := `INSERT INTO chat_history (id, document_id, user_id, message_type, message_content, timestamp) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := database.DB.Exec(query, message.ID, message.DocumentID, message.UserID, message.MessageType, message.MessageContent, message.Timestamp)
	return message, err
}
