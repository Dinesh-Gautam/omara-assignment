package models

import "time"

type ChatMessage struct {
	ID                string    `json:"id"`
	DocumentID        string    `json:"document_id"`
	UserID            string    `json:"user_id"`
	MessageType       string    `json:"message_type"`
	MessageContent    string    `json:"message_content"`
	Timestamp         time.Time `json:"timestamp"`
	AttachedDocuments string    `json:"attached_documents,omitempty"`
}
