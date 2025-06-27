"use client";

import { memo } from "react";
import AnimatedListItem from "@/components/utils/AnimatedListItem";
import { Button } from "@/components/ui/button";
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "@/components/ui/tooltip";
import { downloadDocument } from "@/services/api/documentService";
import { Document } from "@/types/document";
import {
  AlertCircle,
  Download,
  File,
  MessageCircle,
  Trash2,
} from "lucide-react";
import { useRouter } from "next/navigation";
import { PATHS } from "@/constants/config";

/**
 * Props for the DocumentListItem component.
 */
interface DocumentListItemProps {
  /** The document to display. */
  doc: Document;
  /** The index of the item in the list, used for animation. */
  index: number;
  /** A function to remove the document. */
  removeDocument: (documentId: string) => void;
}

/**
 * A component that displays a single document in a list.
 */
const DocumentListItem = memo(
  ({ doc, index, removeDocument }: DocumentListItemProps) => {
    const router = useRouter();

    const handleDelete = async (documentId: string) => {
      if (window.confirm("Are you sure you want to delete this document?")) {
        removeDocument(documentId);
      }
    };

    return (
      <AnimatedListItem key={doc.id} index={index}>
        <div className="flex flex-col rounded-lg border border-gray-200 bg-white p-4 shadow-sm transition-all duration-200 ease-in-out hover:shadow-md md:flex-row md:items-center md:justify-between">
          <div className="flex items-center space-x-4">
            <File className="h-8 w-8 flex-shrink-0 text-gray-400" />
            <div className="flex flex-col">
              <span className="font-semibold text-gray-800 break-all">
                {doc.file_name}
              </span>
              <div className="mt-1 flex flex-wrap items-center gap-2">
                <span className="text-xs text-gray-500">
                  Uploaded at:{" "}
                  {new Date(doc.created_at).toLocaleDateString("en-IN", {
                    day: "2-digit",
                    month: "short",
                    year: "numeric",
                  })}
                </span>
                {doc.status === "processing" && (
                  <span className="rounded-full bg-yellow-100 px-2.5 py-0.5 text-xs font-medium text-yellow-800">
                    Processing
                  </span>
                )}
                {doc.status === "processed" && (
                  <span className="rounded-full bg-green-100 px-2.5 py-0.5 text-xs font-medium text-green-800">
                    Completed
                  </span>
                )}
                {doc.status === "failed" && (
                  <div className="flex items-center space-x-1">
                    <span className="rounded-full bg-red-100 px-2.5 py-0.5 text-xs font-medium text-red-800">
                      Failed
                    </span>
                    <TooltipProvider>
                      <Tooltip>
                        <TooltipTrigger>
                          <AlertCircle className="h-4 w-4 text-red-500" />
                        </TooltipTrigger>
                        <TooltipContent>
                          <p>{doc.processingError}</p>
                        </TooltipContent>
                      </Tooltip>
                    </TooltipProvider>
                  </div>
                )}
              </div>
            </div>
          </div>
          <div className="mt-4 flex items-center justify-end space-x-1 md:mt-0 md:justify-start">
            <Button
              variant="ghost"
              size="icon"
              title="Chat with document"
              onClick={() => router.push(`${PATHS.Chat}/${doc.id}`)}
            >
              <MessageCircle className="h-5 w-5" />
            </Button>
            <Button
              variant="ghost"
              size="icon"
              title="Download document"
              onClick={() => downloadDocument(doc.id, doc.file_name)}
            >
              <Download className="h-5 w-5" />
            </Button>
            <Button
              variant="ghost"
              size="icon"
              title="Delete document"
              className="text-red-500 hover:text-red-700"
              onClick={() => handleDelete(doc.id)}
            >
              <Trash2 className="h-5 w-5" />
            </Button>
          </div>
        </div>
      </AnimatedListItem>
    );
  }
);

export default DocumentListItem;
