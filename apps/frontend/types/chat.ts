/**
 * Represents a single chat message.
 */
export interface ChatMessage {
  /** The unique identifier for the message. */
  id: string;
  /** The ID of the user who sent the message. */
  user_id: string;
  /** The ID of the document this chat is associated with. */
  document_id: string;
  /** The type of message, either from a user or the AI. */
  message_type: "user" | "ai";
  /** The content of the message. */
  message_content: string;
  /** The timestamp when the message was created. */
  timestamp: string;
  /** Optional array of documents attached to the message. */
  attachedDocuments?: { id: string; title: string }[];
}
