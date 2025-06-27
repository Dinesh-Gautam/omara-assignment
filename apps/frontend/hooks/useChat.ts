import { useState, useCallback, useEffect } from "react";
import { flushSync } from "react-dom";
import { ChatMessage } from "../types";
import {
  getChatHistory as apiGetChatHistory,
  postChatMessage as apiPostChatMessage,
} from "../services/api/chatService";
import { v4 as uuidv4 } from "uuid";

/**
 * Custom hook for managing chat functionality.
 * @param documentId - The ID of the document to chat with.
 * @returns An object with chat messages, loading state, error state, and functions to send messages and fetch history.
 */
export const useChat = (documentId: string) => {
  const [messages, setMessages] = useState<ChatMessage[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  /**
   * Fetches the chat history for the current document.
   */
  const fetchHistory = useCallback(async () => {
    if (!documentId) return;
    setIsLoading(true);
    setError(null);
    try {
      const history = await apiGetChatHistory(documentId);
      setMessages(history);
    } catch (err) {
      setError("Failed to fetch chat history. Please try again later.");
      console.error(err);
    } finally {
      setIsLoading(false);
    }
  }, [documentId]);

  useEffect(() => {
    fetchHistory();
  }, [fetchHistory]);

  /**
   * Sends a new message to the chat.
   * @param messageContent - The content of the message to send.
   */
  const sendMessage = async (
    messageContent: string,
    attached_documents: { id: string; title: string }[]
  ) => {
    if (!documentId) return;

    setIsLoading(true);
    const userMessage: ChatMessage = {
      id: uuidv4(),
      document_id: documentId,
      user_id: "current_user",
      message_type: "user",
      message_content: messageContent,
      timestamp: new Date().toISOString(),
      attachedDocuments: attached_documents,
    };
    setMessages((prev) => [...prev, userMessage]);

    const aiMessagePlaceholder: ChatMessage = {
      id: uuidv4(),
      document_id: documentId,
      user_id: "ai",
      message_type: "ai",
      message_content: "",
      timestamp: new Date().toISOString(),
    };
    setMessages((prev) => [...prev, aiMessagePlaceholder]);

    try {
      await apiPostChatMessage(
        documentId,
        messageContent,
        attached_documents.map((d) => d.id),
        (chunk) => {
          flushSync(() => {
            setMessages((prev) =>
              prev.map((msg) =>
                msg.id === aiMessagePlaceholder.id
                  ? { ...msg, message_content: msg.message_content + chunk }
                  : msg
              )
            );
          });
        }
      );
    } catch (err: any) {
      setError(err?.message || "Failed to get a response from the AI.");
      console.error(err);
      setMessages((prev) =>
        prev.filter((msg) => msg.id !== aiMessagePlaceholder.id)
      );
    } finally {
      setIsLoading(false);
    }
  };

  return { messages, isLoading, error, sendMessage, fetchHistory, setMessages };
};
