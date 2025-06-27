/**
 * Represents a document.
 */
export interface Document {
  /** The unique identifier for the document. */
  id: string;
  /** The name of the document file. */
  file_name: string;
  /** The timestamp when the document was uploaded. */
  created_at: string;
  /** The path to the document in Google Cloud Storage. */
  gcs_path: string;
  /** The ID of the user who uploaded the document. */
  user_id: string;
  /** The processing status of the document. */
  status: string;
  /** An optional error message if processing failed. */
  processingError?: string;
}
