"use client";

import { Button } from "@/components/ui/button";
import { Upload } from "lucide-react";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import FileUpload from "./FileUpload";
import { useState } from "react";

const UploadButton = ({
  onUploadComplete,
}: {
  onUploadComplete: (newDoc: any) => void;
}) => {
  const [isOpen, setIsOpen] = useState(false);

  const handleUploadComplete = (newDoc: any) => {
    onUploadComplete(newDoc);
    setIsOpen(false);
  };

  return (
    <Dialog open={isOpen} onOpenChange={setIsOpen}>
      <DialogTrigger asChild>
        <Button>
          <Upload className="mr-2 h-4 w-4" />
          Upload Document
        </Button>
      </DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Upload your PDF</DialogTitle>
        </DialogHeader>
        <FileUpload onUploadComplete={handleUploadComplete} />
      </DialogContent>
    </Dialog>
  );
};

export default UploadButton;
