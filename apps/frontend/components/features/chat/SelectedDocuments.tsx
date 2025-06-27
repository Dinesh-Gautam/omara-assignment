"use client";

import { Document } from "@/types";
import { Button } from "@/components/ui/button";
import { File, X } from "lucide-react";

interface SelectedDocumentsProps {
  selectedDocuments: Document[];
  onRemove: (documentId: string) => void;
}

const SelectedDocuments = ({
  selectedDocuments,
  onRemove,
}: SelectedDocumentsProps) => {
  if (selectedDocuments.length === 0) {
    return null;
  }

  return (
    <div className="p-2 pt-0 border-gray-200">
      <div className="flex flex-wrap gap-2">
        {selectedDocuments.map((doc) => (
          <div
            key={doc.id}
            className="flex items-center bg-gray-200 rounded-full px-3 py-1 text-sm"
          >
            <File className="h-4 w-4 mr-2" />
            <span>{doc.file_name}</span>
            <Button
              variant="ghost"
              size="sm"
              className="ml-2 h-6 w-6 p-0 rounded-full"
              onClick={() => onRemove(doc.id)}
            >
              <X className="h-4 w-4" />
            </Button>
          </div>
        ))}
      </div>
    </div>
  );
};

export default SelectedDocuments;
