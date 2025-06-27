import { auth } from "@/lib/firebase";
import axios from "../../lib/axios";
import { ChatMessage } from "../../types";

/**
 * Fetches the chat history for a specific document.
 * @param documentId - The ID of the document.
 * @returns A promise that resolves to an array of chat messages.
 */
export const getChatHistory = async (
  documentId: string
): Promise<ChatMessage[]> => {
  const response = await axios.get(`/api/chat/${documentId}`);
  // The backend returns attached_documents as a JSON string.
  // We need to parse it on the frontend.
  return response.data.map((message: any) => {
    if (
      message.attached_documents &&
      typeof message.attached_documents === "string"
    ) {
      try {
        const parsedOnce = JSON.parse(message.attached_documents);
        return {
          ...message,
          attachedDocuments: parsedOnce,
        };
      } catch (e) {
        console.error(
          "Failed to parse attached_documents:",
          e,
          "Raw data:",
          message.attached_documents
        );
        return { ...message, attachedDocuments: [] };
      }
    }
    return { ...message, attachedDocuments: message.attached_documents || [] };
  });
};

/**
 * Posts a chat message and handles streaming response.
 * @param documentId - The ID of the document.
 * @param message - The message to post.
 * @param onChunk - A callback function to handle each chunk of the response.
 * @returns A promise that resolves when the stream is complete.
 */
export const postChatMessage = async (
  documentId: string,
  message: string,
  attached_documents: string[],
  onChunk: (chunk: string) => void
): Promise<void> => {
  const user = auth.currentUser;
  if (!user) {
    throw new Error("User not authenticated");
  }
  const token = await user.getIdToken();
  const authString = `Bearer ${token}`;

  const response = await fetch(process.env.NEXT_PUBLIC_API_URL + "/api/chat", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      Accept: "text/event-stream",
      Authorization: authString,
    },
    body: JSON.stringify({
      document_id: documentId,
      message: message,
      attached_documents: attached_documents,
    }),
  });

  if (!response.ok) {
    // Handle non-streaming errors (e.g., 4xx, 5xx)
    const errorData = await response.json();
    throw new Error(errorData.error || "An unexpected error occurred");
  }

  if (!response.body) {
    throw new Error("No response body");
  }

  const reader = response.body.getReader();
  const decoder = new TextDecoder();
  let buffer = "";

  try {
    while (true) {
      const { done, value } = await reader.read();
      if (done) break;

      const chunk = decoder.decode(value, { stream: true });
      buffer += chunk;

      const lines = buffer.split("\n");
      buffer = lines.pop() || "";

      for (const line of lines) {
        if (line.startsWith("data: ")) {
          const data = line.substring(6);

          if (data.trim() === "[DONE]") {
            return;
          }

          try {
            const json = JSON.parse(data);
            if (json.error) {
              // Handle structured error from the stream
              throw new Error(json.error);
            }
            if (json.token) {
              onChunk(json.token);
            }
          } catch (e) {
            console.error("Failed to parse stream chunk:", e, "Data:", data);
          }
        }
      }
    }
  } finally {
    reader.releaseLock();
  }
};
