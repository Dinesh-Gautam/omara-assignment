"use client";

import ChatInterface from "@/components/features/chat/ChatInterface";
import { useParams } from "next/navigation";

export default function ChatPage() {
  const { documentId } = useParams();

  return <ChatInterface documentId={documentId as string} />;
}
