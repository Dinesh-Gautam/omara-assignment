package processor

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

// ExtractText handles non-streaming text extraction for simple file types.
func ExtractText(file io.Reader, fileType string) (string, error) {
	switch fileType {
	case "text/plain":
		return extractTextFromTXT(file)
	default:
		return "", fmt.Errorf("unsupported file type for non-streaming extraction: %s", fileType)
	}
}

// ProcessPDFChunks extracts text from a PDF, chunks it, and processes each chunk via a callback.
// It uses the `pdftotext` command-line tool.
// NOTE: This function requires the `poppler-utils` package (which provides `pdftotext`)
// to be installed on the system running the backend.
func ProcessPDFChunks(file io.Reader, chunkSize int, overlap int, processChunk func(chunk string, chunkIndex int) error) error {
	// Create a temporary file for the uploaded PDF
	inputFile, err := ioutil.TempFile("", "upload-*.pdf")
	if err != nil {
		return fmt.Errorf("failed to create temp input file: %w", err)
	}
	defer os.Remove(inputFile.Name())

	// Copy the uploaded file content to the temporary file
	if _, err := io.Copy(inputFile, file); err != nil {
		return fmt.Errorf("failed to copy to temp file: %w", err)
	}
	inputFile.Close()

	// Create a temporary file for the text output
	outputFile, err := ioutil.TempFile("", "output-*.txt")
	if err != nil {
		return fmt.Errorf("failed to create temp output file: %w", err)
	}
	defer os.Remove(outputFile.Name())
	outputFile.Close() // Close the file so pdftotext can write to it

	// Execute the pdftotext command
	// The -layout flag helps preserve the document's structure.
	cmd := exec.Command("pdftotext", "-layout", inputFile.Name(), outputFile.Name())
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run pdftotext command: %w. Ensure poppler-utils is installed", err)
	}

	// Read the entire text file content
	textContent, err := ioutil.ReadFile(outputFile.Name())
	if err != nil {
		return fmt.Errorf("failed to read text output file: %w", err)
	}

	// Use the existing rune-safe ChunkText function
	chunks := ChunkText(string(textContent), chunkSize, overlap)

	// Process each chunk
	for i, chunk := range chunks {
		if err := processChunk(chunk, i); err != nil {
			// If the callback returns an error, abort the processing and return the error.
			return fmt.Errorf("failed to process chunk %d: %w", i, err)
		}
	}

	return nil
}

func extractTextFromTXT(file io.Reader) (string, error) {
	var text strings.Builder
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text.WriteString(scanner.Text())
		text.WriteString("\n")
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}
	return text.String(), nil
}

func ChunkText(text string, chunkSize int, overlap int) []string {
	var chunks []string
	runes := []rune(text)
	if len(runes) == 0 {
		return chunks
	}

	for i := 0; i < len(runes); i += chunkSize - overlap {
		end := i + chunkSize
		if end > len(runes) {
			end = len(runes)
		}
		chunks = append(chunks, string(runes[i:end]))
		if end == len(runes) {
			break
		}
	}
	return chunks
}
