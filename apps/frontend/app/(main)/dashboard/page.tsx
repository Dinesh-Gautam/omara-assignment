"use client";

import DocumentList from "@/components/features/dashboard/DocumentList";
import UploadButton from "@/components/features/dashboard/UploadButton";
import { Button } from "@/components/ui/button";
import { useAuth } from "@/context/AuthContext";
import { useDocuments } from "@/hooks/useDocuments";
import { auth } from "@/lib/firebase";

export default function Dashboard() {
  const { user } = useAuth();

  const {
    documents,
    isLoading: docsLoading,
    removeDocument,
    addDocument,
  } = useDocuments();

  return (
    user && (
      <div className="container mx-auto p-4 md:p-8">
        <div className="flex justify-between items-center mb-8">
          <div>
            <h1 className="text-3xl font-bold">Dashboard</h1>
            <p className="text-muted-foreground">
              Welcome, {user.isAnonymous ? "Guest" : user.email}
            </p>
          </div>
          <Button onClick={() => auth.signOut()} variant="outline">
            Sign Out
          </Button>
        </div>

        <div className="flex flex-col items-start gap-4">
          <UploadButton onUploadComplete={addDocument} />
          <DocumentList
            documents={documents}
            isLoading={docsLoading}
            removeDocument={removeDocument}
          />
        </div>
      </div>
    )
  );
}
