import axios from "../../lib/axios";
import { Document } from "../../types";

/**
 * Uploads a document.
 * @param file - The file to upload.
 * @param onUploadProgress - A callback to track upload progress.
 * @returns A promise that resolves to the uploaded document.
 */
export const uploadDocument = async (
  file: File,
  onUploadProgress: (progress: number) => void
): Promise<Document> => {
  const formData = new FormData();
  formData.append("file", file);

  const response = await axios.post("/api/documents/upload", formData, {
    headers: {
      "Content-Type": "multipart/form-data",
    },
    onUploadProgress: (progressEvent) => {
      const percentCompleted = Math.round(
        (progressEvent.loaded * 100) / (progressEvent.total ?? 1)
      );
      onUploadProgress(percentCompleted);
    },
  });

  return response.data;
};

/**
 * Checks the status of a document.
 * @param documentId - The ID of the document.
 * @returns A promise that resolves to the document with its current status.
 */
export const checkDocumentStatus = async (
  documentId: string
): Promise<Document> => {
  const response = await axios.get(`/api/documents/${documentId}/status`);
  return response.data;
};

/**
 * Fetches all documents for the user.
 * @returns A promise that resolves to an array of documents.
 */
export const getDocuments = async (): Promise<Document[]> => {
  const response = await axios.get("/api/documents");
  return response.data;
};

/**
 * Downloads a document.
 * @param documentId - The ID of the document to download.
 * @param fileName - The name to use for the downloaded file.
 */
export const downloadDocument = async (
  documentId: string,
  fileName: string
): Promise<void> => {
  const response = await axios.get(`/api/documents/download/${documentId}`, {
    responseType: "blob",
  });

  const url = window.URL.createObjectURL(new Blob([response.data]));
  const a = document.createElement("a");
  a.href = url;
  a.download = fileName;
  document.body.appendChild(a);
  a.click();
  a.remove();
  window.URL.revokeObjectURL(url);
};

/**
 * Deletes a document.
 * @param documentId - The ID of the document to delete.
 */
export const deleteDocument = async (documentId: string): Promise<void> => {
  await axios.delete(`/api/documents/${documentId}`);
};
