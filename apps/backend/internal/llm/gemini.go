package llm

import (
	"context"
	"flag"
	"fmt"
	"strategic-insight-analyst/backend/config"
	"strings"

	"google.golang.org/genai"
)

var model = flag.String("model", "gemini-2.0-flash", "the model name, e.g. gemini-2.0-flash")

// GetEmbedding generates an embedding for the given text using the Gemini API.
func GetEmbedding(text string) ([]float32, error) {
	ctx := context.Background()
	apiKey := config.AppConfig.GeminiAPIKey
	if apiKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY not set in config")
	}

	client, err := genai.NewClient(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	contents := []*genai.Content{
		genai.NewContentFromText(text, genai.RoleUser),
	}

	result, err := client.Models.EmbedContent(ctx, "text-embedding-004", contents, &genai.EmbedContentConfig{OutputDimensionality: genai.Ptr[int32](768)})

	if err != nil {
		return nil, fmt.Errorf("failed to generate embedding: %w", err)
	}

	if result != nil && len(result.Embeddings) > 0 && result.Embeddings[0] != nil {
		return result.Embeddings[0].Values, nil
	}

	return nil, fmt.Errorf("no embedding found in response")
}

func CallGeminiStream(query string, contextText string, history []*genai.Content, hasAttachedDocs bool, streamChan chan<- string) (string, error) {
	defer close(streamChan)
	ctx := context.Background()
	apiKey := config.AppConfig.GeminiAPIKey
	if apiKey == "" {
		return "", fmt.Errorf("GEMINI_API_KEY not set in config")
	}

	systemInstruction := `You are a sophisticated AI assistant specializing in strategic analysis. Your primary function is to deliver precise, insightful, and concise answers based *exclusively* on the provided document context.

Key Instructions:
1. **Strict Context Adherence:** Base your analysis *only* on the text within the '--- Document Context ---' or '<document_context>' section. Do not use any external knowledge or make assumptions.
2. **Acknowledge Limitations:** If the information required to answer the query is not present in the provided context, you *must* explicitly state that the information is not available. Do not attempt to invent or infer information.
3. **Clear & Concise Output:** Present your analysis in a clear and easily digestible format. The user's query may specify a desired format (e.g., a bulleted list).`

	config := &genai.GenerateContentConfig{
		Temperature:       genai.Ptr[float32](0.5),
		SystemInstruction: genai.NewContentFromText(systemInstruction, genai.RoleUser),
	}

	client, err := genai.NewClient(ctx, nil)
	if err != nil {
		return "", err
	}

	chat, err := client.Chats.Create(ctx, *model, config, history)
	if err != nil {
		return "", err
	}

	var prompt string
	if hasAttachedDocs {
		prompt = fmt.Sprintf(`Based on the context provided in the following document(s), answer the user's question.

<document_context>
%s
</document_context>

User Question: %s`, contextText, query)
	} else {
		prompt = fmt.Sprintf(`Based *only* on the provided document context, analyze and answer the following strategic query.

--- Document Context ---
%s
--- End Context ---

User Query: %s

Provide a clear, concise, and analytical response. If the query suggests a format (e.g., 'Provide a bulleted list of 3 key insights'), adhere to it. If the necessary information is not in the context, state that clearly.`, contextText, query)
	}

	stream := chat.SendMessageStream(ctx, genai.Part{Text: prompt})
	var fullResponse strings.Builder
	for chunk, _ := range stream {

		part := chunk.Candidates[0].Content.Parts[0]
		streamChan <- part.Text
		fullResponse.WriteString(part.Text)

	}

	return fullResponse.String(), nil
}
