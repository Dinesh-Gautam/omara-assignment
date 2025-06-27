"use client";

import { Button } from "@/components/ui/button";
import { Progress } from "@/components/ui/progress";
import { Document } from "@/types/document";
import { useFileUpload } from "@/hooks/useFileUpload";
import { File as FileIcon, UploadCloud, X } from "lucide-react";
import { useCallback, useState } from "react";
import { useDropzone } from "react-dropzone";

interface FileDropzoneProps {
  onUploadComplete: (document: Document) => void;
}

const FileDropzone = ({ onUploadComplete }: FileDropzoneProps) => {
  const [file, setFile] = useState<File | null>(null);
  const { uploadProgress, isUploading, error, handleUpload, setError } =
    useFileUpload({
      onUploadComplete: (document) => {
        onUploadComplete(document);
        setFile(null);
      },
    });

  const onDrop = useCallback(
    (acceptedFiles: File[], rejectedFiles: any[]) => {
      setError(null);
      if (rejectedFiles.length > 0) {
        const errorMessage = rejectedFiles[0].errors[0].message;
        setError(errorMessage);
        return;
      }

      if (acceptedFiles.length > 0) {
        setFile(acceptedFiles[0]);
      }
    },
    [setError]
  );

  const { getRootProps, getInputProps, isDragActive } = useDropzone({
    onDrop,
    accept: {
      "application/pdf": [".pdf"],
      "text/plain": [".txt"],
    },
    maxSize: 10 * 1024 * 1024, // 10MB
    multiple: false,
  });

  const onRemoveFile = () => {
    setFile(null);
    setError(null);
  };

  const onUpload = () => {
    if (file) {
      handleUpload(file);
    }
  };

  return (
    <>
      {file ? (
        <div className="flex flex-col items-center justify-center space-y-4">
          <div className="flex items-center space-x-4">
            <FileIcon className="h-10 w-10 text-gray-500" />
            <div className="text-left">
              <p className="font-medium">{file.name}</p>
              <p className="text-sm text-gray-500">
                {(file.size / 1024 / 1024).toFixed(2)} MB
              </p>
            </div>
            <Button
              variant="ghost"
              size="icon"
              onClick={onRemoveFile}
              disabled={isUploading}
            >
              <X className="h-5 w-5" />
            </Button>
          </div>
          {isUploading && (
            <Progress value={uploadProgress} className="w-full" />
          )}
          <Button onClick={onUpload} disabled={isUploading}>
            {isUploading
              ? `Uploading... ${uploadProgress.toFixed(0)}%`
              : "Upload"}
          </Button>
        </div>
      ) : (
        <div
          {...getRootProps()}
          className={`flex flex-col items-center justify-center space-y-4 rounded-lg border-2 border-dashed p-12 text-center transition-colors ${
            isDragActive ? "border-primary bg-primary/10" : "border-gray-300"
          }`}
        >
          <input {...getInputProps()} />
          <UploadCloud className="h-12 w-12 text-gray-400" />
          <p className="text-lg font-semibold">
            Drag & drop a file here, or click to select one
          </p>
          <p className="text-sm text-gray-500">PDF or TXT files, up to 10MB</p>
        </div>
      )}
      {error && (
        <p className="mt-2 text-center text-sm text-red-500">{error}</p>
      )}
    </>
  );
};

export default FileDropzone;
