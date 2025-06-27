"use client";

import { memo } from "react";
import FileDropzone from "@/components/utils/FileDropzone";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Document } from "@/types/document";

/**
 * Props for the FileUpload component.
 */
interface FileUploadProps {
  /** A callback function that is called when a file upload is complete. */
  onUploadComplete: (document: Document) => void;
}

/**
 * A component for uploading files.
 */
const FileUpload = memo(({ onUploadComplete }: FileUploadProps) => {
  return (
    <Card>
      <CardHeader>
        <CardTitle>Upload Document</CardTitle>
        <CardDescription>
          Drag and drop a file or click to select.
        </CardDescription>
      </CardHeader>
      <CardContent>
        <FileDropzone onUploadComplete={onUploadComplete} />
      </CardContent>
    </Card>
  );
});

export default FileUpload;
