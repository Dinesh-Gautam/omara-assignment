import { useState, useEffect, useCallback, useRef } from "react";
import { toast } from "sonner";
import { Document } from "../types";
import {
  getDocuments as apiGetDocuments,
  deleteDocument as apiDeleteDocument,
  checkDocumentStatus,
} from "../services/api/documentService";

/**
 * Custom hook for managing documents.
 * @returns An object with the list of documents, loading state, error state, and functions to fetch, add, and remove documents.
 */
export const useDocuments = () => {
  const [documents, setDocuments] = useState<Document[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  /**
   * Fetches the list of documents from the server.
   */
  const fetchDocuments = useCallback(async () => {
    setIsLoading(true);
    setError(null);
    try {
      const data = await apiGetDocuments();
      setDocuments(data);
    } catch (err) {
      setError("Failed to fetch documents");
      console.error(err);
    } finally {
      setIsLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchDocuments();
  }, [fetchDocuments]);

  const pollingRef = useRef<Set<string>>(new Set());

  useEffect(() => {
    const poll = async (doc: Document) => {
      if (!pollingRef.current.has(doc.id)) {
        return;
      }

      try {
        const updatedDocument = await checkDocumentStatus(doc.id);

        if (updatedDocument.status !== "processing") {
          pollingRef.current.delete(doc.id);
          setDocuments((prevDocs) =>
            prevDocs.map((d) =>
              d.id === updatedDocument.id ? updatedDocument : d
            )
          );
          if (updatedDocument.status === "failed") {
            toast.error(
              `Processing failed for document: ${updatedDocument.file_name}`
            );
          }
        } else {
          setDocuments((prevDocs) =>
            prevDocs.map((d) =>
              d.id === updatedDocument.id ? updatedDocument : d
            )
          );
          setTimeout(() => poll(updatedDocument), 1000);
        }
      } catch (err) {
        console.error(`Error polling status for document ${doc.id}:`, err);
        toast.error(`Could not update status for ${doc.file_name}`);

        pollingRef.current.delete(doc.id);

        setDocuments((prevDocs) =>
          prevDocs.map((d) =>
            d.id === doc.id ? { ...d, status: "failed" } : d
          )
        );
      }
    };

    documents.forEach((doc) => {
      if (doc.status === "processing" && !pollingRef.current.has(doc.id)) {
        pollingRef.current.add(doc.id);
        poll(doc);
      }
    });
  }, [documents]);

  /**
   * Adds a new document to the local state.
   * @param document - The document to add.
   */
  const addDocument = (document: Document) => {
    setDocuments((prevDocs) => [document, ...prevDocs]);
  };

  /**
   * Removes a document from the local state and deletes it from the server.
   * @param documentId - The ID of the document to remove.
   */
  const removeDocument = async (documentId: string) => {
    setDocuments((prevDocs) => prevDocs.filter((doc) => doc.id !== documentId));
    try {
      await apiDeleteDocument(documentId);
    } catch (err) {
      setError("Failed to delete document");
      console.error(err);
      // revert state if API call fails
      fetchDocuments();
    }
  };

  return {
    documents,
    isLoading,
    error,
    fetchDocuments,
    addDocument,
    removeDocument,
  };
};
