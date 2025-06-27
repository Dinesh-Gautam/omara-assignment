"use client";

import { useState, useEffect } from "react";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { useDocuments } from "@/hooks/useDocuments";
import { Document } from "@/types";
import { Input } from "@/components/ui/input";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Checkbox } from "@/components/ui/checkbox";

interface DocumentSelectionModalProps {
  isOpen: boolean;
  onClose: () => void;
  onConfirm: (selectedIds: string[]) => void;
  initialSelectedIds: string[];
}

const DocumentSelectionModal = ({
  isOpen,
  onClose,
  onConfirm,
  initialSelectedIds,
}: DocumentSelectionModalProps) => {
  const { documents, isLoading, error } = useDocuments();
  const [selectedIds, setSelectedIds] = useState<string[]>(initialSelectedIds);
  const [searchTerm, setSearchTerm] = useState("");

  useEffect(() => {
    setSelectedIds(initialSelectedIds);
  }, [initialSelectedIds]);

  const handleToggleSelection = (id: string) => {
    setSelectedIds((prev) =>
      prev.includes(id) ? prev.filter((docId) => docId !== id) : [...prev, id]
    );
  };

  const filteredDocuments = documents.filter((doc) =>
    doc.file_name.toLowerCase().includes(searchTerm.toLowerCase())
  );

  const handleConfirm = () => {
    onConfirm(selectedIds);
    onClose();
  };

  return (
    <Dialog open={isOpen} onOpenChange={onClose}>
      <DialogContent className="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>Select Documents for Context</DialogTitle>
        </DialogHeader>
        <div className="py-4">
          <Input
            placeholder="Search documents..."
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
            className="mb-4"
          />
          <ScrollArea className="h-72 w-full rounded-md border">
            <div className="p-4">
              {isLoading ? (
                <p>Loading documents...</p>
              ) : error ? (
                <p className="text-red-500">{error}</p>
              ) : (
                filteredDocuments.map((doc) => (
                  <div
                    key={doc.id}
                    className="flex items-center space-x-2 mb-2"
                  >
                    <Checkbox
                      id={doc.id}
                      checked={selectedIds.includes(doc.id)}
                      onCheckedChange={() => handleToggleSelection(doc.id)}
                    />
                    <label
                      htmlFor={doc.id}
                      className="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70"
                    >
                      {doc.file_name}
                    </label>
                  </div>
                ))
              )}
            </div>
          </ScrollArea>
        </div>
        <DialogFooter>
          <Button variant="outline" onClick={onClose}>
            Cancel
          </Button>
          <Button onClick={handleConfirm}>Confirm</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
};

export default DocumentSelectionModal;
