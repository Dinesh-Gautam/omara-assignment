"use client";

import { memo } from "react";
import AnimatedList from "@/components/utils/AnimatedList";
import DocumentListSkeleton from "@/components/utils/DocumentListSkeleton";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Document } from "@/types/document";
import DocumentListItem from "./DocumentListItem";

/**
 * Props for the DocumentList component.
 */
interface DocumentListProps {
  /** The list of documents to display. */
  documents: Document[];
  /** A flag indicating if the documents are being loaded. */
  isLoading: boolean;
  /** A function to remove a document. */
  removeDocument: (documentId: string) => void;
}

/**
 * A component that displays a list of documents.
 */
const DocumentList = memo(
  ({ documents, isLoading, removeDocument }: DocumentListProps) => {
    return (
      <Card className="w-full">
        <CardHeader>
          <CardTitle>My Documents</CardTitle>
          <CardDescription>View your uploaded documents.</CardDescription>
        </CardHeader>
        <CardContent>
          {isLoading ? (
            <DocumentListSkeleton />
          ) : documents.length > 0 ? (
            <AnimatedList>
              {documents.map((doc, index) => (
                <DocumentListItem
                  key={doc.id}
                  doc={doc}
                  index={index}
                  removeDocument={removeDocument}
                />
              ))}
            </AnimatedList>
          ) : (
            <p>No documents uploaded yet.</p>
          )}
        </CardContent>
      </Card>
    );
  }
);

export default DocumentList;
