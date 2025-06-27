"use client";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import ChatSkeleton from "@/components/utils/ChatSkeleton";
import { useChat } from "@/hooks/useChat";
import { useDocuments } from "@/hooks/useDocuments";
import { Bot, File, Paperclip, Send, User } from "lucide-react";
import { useEffect, useRef, useState } from "react";
import ReactMarkdown from "react-markdown";
import DocumentSelectionModal from "./DocumentSelectionModal";
import SelectedDocuments from "./SelectedDocuments";

/**
 * Props for the ChatInterface component.
 */
interface ChatInterfaceProps {
  /** The ID of the document to chat with. */
  documentId: string;
}

/**
 * A component that provides a chat interface for a specific document.
 */
const ChatInterface = ({ documentId }: ChatInterfaceProps) => {
  const { messages, isLoading, error, sendMessage, setMessages } =
    useChat(documentId);
  const [newMessage, setNewMessage] = useState("");
  const messagesEndRef = useRef<HTMLDivElement>(null);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [selectedDocIds, setSelectedDocIds] = useState<string[]>([]);
  const { documents } = useDocuments();

  const selectedDocuments = documents.filter((doc) =>
    selectedDocIds.includes(doc.id)
  );

  const handleRemoveDocument = (docId: string) => {
    setSelectedDocIds((prev) => prev.filter((id) => id !== docId));
  };

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
  };

  useEffect(() => {
    scrollToBottom();
  }, [messages]);

  useEffect(() => {
    if (error) {
      const errorId = `error-${Date.now()}`;
      setMessages((prev) => [
        ...prev,
        {
          id: errorId,
          document_id: documentId,
          user_id: "ai",
          message_type: "ai",
          message_content: `Error: ${error}`,
          timestamp: new Date().toISOString(),
        },
      ]);
    }
  }, [error, documentId, setMessages]);

  const handleSendMessage = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!newMessage.trim()) return;
    sendMessage(
      newMessage,
      selectedDocuments.map((d) => ({ id: d.id, title: d.file_name }))
    );
    setNewMessage("");
    setSelectedDocIds([]);
  };

  const sampleQuestions = [
    "What is the main topic of this document?",
    "Summarize the key points in 3 bullets.",
    "What are the main conclusions?",
  ];

  const handleSampleMessageClick = (question: string) => {
    sendMessage(
      question,
      selectedDocuments.map((d) => ({ id: d.id, title: d.file_name }))
    );
  };

  if (isLoading && messages.length === 0) {
    return <ChatSkeleton />;
  }

  return (
    <div className="flex h-screen flex-col bg-gray-50">
      <header className="flex h-16 flex-shrink-0 items-center justify-between border-b border-gray-200 bg-white px-4 md:px-6">
        <h1 className="text-lg font-semibold">Chat with Document</h1>
      </header>
      <main className="flex-1 overflow-y-auto p-4 md:p-6">
        <div className="space-y-6">
          {messages.length === 0 ? (
            <div className="flex flex-col items-center justify-center h-full text-center">
              <h2 className="text-xl font-semibold text-gray-600">
                No messages yet
              </h2>
              <p className="text-gray-500">
                Start the conversation by asking a question.
              </p>
              <div className="mt-4 flex flex-wrap justify-center gap-2">
                {sampleQuestions.map((q) => (
                  <Button
                    key={q}
                    variant="outline"
                    onClick={() => handleSampleMessageClick(q)}
                  >
                    {q}
                  </Button>
                ))}
              </div>
            </div>
          ) : (
            messages.map((msg) => (
              <div
                key={msg.id}
                className={`flex items-start gap-4 ${
                  msg.message_type === "user" ? "justify-end" : ""
                }`}
              >
                {msg.message_type === "ai" && (
                  <div className="flex h-10 w-10 flex-shrink-0 items-center justify-center rounded-full bg-gray-200">
                    <Bot className="h-6 w-6 text-gray-600" />
                  </div>
                )}
                <div
                  className={`max-w-2xl rounded-lg px-4 py-3 ${
                    msg.message_type === "user"
                      ? "bg-blue-500 text-white"
                      : "bg-white text-gray-800"
                  }`}
                >
                  <div
                    className={
                      msg.message_type === "ai"
                        ? "prose prose-sm max-w-none"
                        : ""
                    }
                  >
                    <ReactMarkdown>{msg.message_content}</ReactMarkdown>
                    {isLoading &&
                      msg.message_type === "ai" &&
                      msg.message_content.length === 0 && (
                        <span className="animate-pulse">...</span>
                      )}
                  </div>
                  {msg.attachedDocuments &&
                    msg.attachedDocuments.length > 0 && (
                      <div className={`mt-1`}>
                        <div className="flex flex-wrap gap-2">
                          {msg.attachedDocuments.map((doc) => (
                            <div
                              key={doc.id}
                              className={`flex items-center rounded-sm px-2 py-1 text-xs bg-primary/20`}
                            >
                              <File className="h-3 w-3 mr-1" />
                              <span>{doc.title}</span>
                            </div>
                          ))}
                        </div>
                      </div>
                    )}
                  <div
                    className={`mt-1 text-xs ${
                      msg.message_type === "user"
                        ? "text-blue-200"
                        : "text-gray-500"
                    }`}
                  >
                    <span className="text-xs">
                      {new Date(msg.timestamp).toLocaleTimeString("en-IN", {
                        hour12: true,
                        hour: "numeric",
                        minute: "numeric",
                      })}
                    </span>
                  </div>
                </div>
                {msg.message_type === "user" && (
                  <div className="flex h-10 w-10 flex-shrink-0 items-center justify-center rounded-full bg-blue-500">
                    <User className="h-6 w-6 text-white" />
                  </div>
                )}
              </div>
            ))
          )}
          <div ref={messagesEndRef} />
        </div>
      </main>
      <footer className="border-t border-gray-200 bg-white p-4 md:p-6">
        <SelectedDocuments
          selectedDocuments={selectedDocuments}
          onRemove={handleRemoveDocument}
        />
        <form onSubmit={handleSendMessage} className="flex items-center gap-2">
          <Button
            type="button"
            variant="outline"
            size="icon"
            onClick={() => setIsModalOpen(true)}
          >
            <Paperclip className="h-4 w-4" />
          </Button>
          <Input
            value={newMessage}
            onChange={(e) => setNewMessage(e.target.value)}
            placeholder="Type your message..."
            className="flex-1"
          />
          <Button type="submit" disabled={isLoading}>
            <Send className="mr-2 h-4 w-4" />
            Send
          </Button>
        </form>
        <DocumentSelectionModal
          isOpen={isModalOpen}
          onClose={() => setIsModalOpen(false)}
          onConfirm={setSelectedDocIds}
          initialSelectedIds={selectedDocIds}
        />
      </footer>
    </div>
  );
};

export default ChatInterface;
