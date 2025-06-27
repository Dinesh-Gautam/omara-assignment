import { uploadDocument } from "@/services/api/documentService";
import { Document } from "@/types/document";
import { AxiosError } from "axios";
import { useState } from "react";

/**
 * Props for the useFileUpload hook.
 */
interface UseFileUploadProps {
  /**
   * Callback function to be executed when the upload is complete.
   * @param document - The uploaded document.
   */
  onUploadComplete: (document: Document) => void;
}

/**
 * Custom hook for handling file uploads.
 * @param onUploadComplete - Callback function to be executed when the upload is complete.
 * @returns An object containing upload progress, status, error, and upload handler.
 */
export const useFileUpload = ({ onUploadComplete }: UseFileUploadProps) => {
  const [uploadProgress, setUploadProgress] = useState(0);
  const [isUploading, setIsUploading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  /**
   * Handles the file upload process.
   * @param file - The file to be uploaded.
   */
  const handleUpload = async (file: File) => {
    if (!file) return;

    setIsUploading(true);
    setError(null);
    setUploadProgress(0);

    try {
      const uploadedDocument = await uploadDocument(file, setUploadProgress);

      setTimeout(() => {
        onUploadComplete(uploadedDocument);
        setIsUploading(false);
      }, 1000);
    } catch (error) {
      if (error instanceof AxiosError) {
        const errorMessage =
          error.response?.data?.error || "Failed to upload file.";
        setError(errorMessage);
      } else {
        setError("An unexpected error occurred.");
      }
      setIsUploading(false);
    }
  };

  return {
    uploadProgress,
    isUploading,
    error,
    handleUpload,
    setError,
  };
};
