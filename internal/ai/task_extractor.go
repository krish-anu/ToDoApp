package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

var (
	ErrMissingAPIKey     = errors.New("openai api key is not configured")
	ErrMalformedResponse = errors.New("openai response missing structured output")
)

type TaskExtraction struct {
	Title    string   `json:"title"`
	DueAt    *string  `json:"due_at"`
	Priority string   `json:"priority"`
	Tags     []string `json:"tags"`
	RemindAt *string  `json:"remind_at"`
}

type TaskExtractor interface {
	ExtractTaskFromText(ctx context.Context, message string, timezone string) (TaskExtraction, error)
}

type OpenAITaskExtractor struct {
	apiKey     string
	model      string
	baseURL    string
	httpClient *http.Client
}

func NewOpenAITaskExtractor(apiKey string, model string, baseURL string) *OpenAITaskExtractor {
	resolvedModel := strings.TrimSpace(model)
	if resolvedModel == "" {
		resolvedModel = "gpt-4.1-mini"
	}

	resolvedBaseURL := strings.TrimSuffix(strings.TrimSpace(baseURL), "/")
	if resolvedBaseURL == "" {
		resolvedBaseURL = "https://api.openai.com"
	}

	return &OpenAITaskExtractor{
		apiKey:  strings.TrimSpace(apiKey),
		model:   resolvedModel,
		baseURL: resolvedBaseURL,
		httpClient: &http.Client{
			Timeout: 25 * time.Second,
		},
	}
}

func (c *OpenAITaskExtractor) ExtractTaskFromText(ctx context.Context, message string, timezone string) (TaskExtraction, error) {
	if c.apiKey == "" {
		return TaskExtraction{}, ErrMissingAPIKey
	}

	trimmedMessage := strings.TrimSpace(message)
	if trimmedMessage == "" {
		return TaskExtraction{}, errors.New("message cannot be empty")
	}

	tz := strings.TrimSpace(timezone)
	if tz == "" {
		tz = "UTC"
	}

	requestPayload := responsesRequest{
		Model: c.model,
		Input: []responsesInputMessage{
			{
				Role: "system",
				Content: []responsesInputContent{
					{Type: "input_text", Text: taskExtractionSystemPrompt},
				},
			},
			{
				Role: "user",
				Content: []responsesInputContent{
					{Type: "input_text", Text: buildExtractionInput(trimmedMessage, tz)},
				},
			},
		},
		Text: responsesTextConfig{
			Format: responsesTextFormat{
				Type:   "json_schema",
				Name:   "todo_task_extraction",
				Schema: taskExtractionJSONSchema(),
				Strict: true,
			},
		},
	}

	requestBody, err := json.Marshal(requestPayload)
	if err != nil {
		return TaskExtraction{}, fmt.Errorf("failed to marshal openai request: %w", err)
	}

	httpRequest, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/v1/responses", bytes.NewReader(requestBody))
	if err != nil {
		return TaskExtraction{}, fmt.Errorf("failed to create openai request: %w", err)
	}
	httpRequest.Header.Set("Authorization", "Bearer "+c.apiKey)
	httpRequest.Header.Set("Content-Type", "application/json")

	httpResponse, err := c.httpClient.Do(httpRequest)
	if err != nil {
		return TaskExtraction{}, fmt.Errorf("failed to call openai: %w", err)
	}
	defer httpResponse.Body.Close()

	responseBody, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		return TaskExtraction{}, fmt.Errorf("failed to read openai response: %w", err)
	}

	if httpResponse.StatusCode >= http.StatusBadRequest {
		if apiErrMessage := parseOpenAIErrorMessage(responseBody); apiErrMessage != "" {
			return TaskExtraction{}, fmt.Errorf("openai error (%d): %s", httpResponse.StatusCode, apiErrMessage)
		}
		return TaskExtraction{}, fmt.Errorf("openai returned status %d: %s", httpResponse.StatusCode, truncateForError(string(responseBody), 300))
	}

	var parsedResponse responsesResponse
	if err := json.Unmarshal(responseBody, &parsedResponse); err != nil {
		return TaskExtraction{}, fmt.Errorf("failed to decode openai response: %w", err)
	}

	if parsedResponse.Error != nil && parsedResponse.Error.Message != "" {
		return TaskExtraction{}, fmt.Errorf("openai error: %s", parsedResponse.Error.Message)
	}

	outputText := strings.TrimSpace(parsedResponse.OutputText)
	if outputText == "" {
		outputText = strings.TrimSpace(readOutputText(parsedResponse.Output))
	}
	if outputText == "" {
		return TaskExtraction{}, ErrMalformedResponse
	}

	var extracted TaskExtraction
	if err := json.Unmarshal([]byte(outputText), &extracted); err != nil {
		return TaskExtraction{}, fmt.Errorf("failed to parse extracted task json: %w", err)
	}

	return normalizeExtraction(extracted), nil
}

func buildExtractionInput(message string, timezone string) string {
	return fmt.Sprintf("Current datetime (UTC): %s\nUser timezone: %s\nUser message: %s", time.Now().UTC().Format(time.RFC3339), timezone, message)
}

func normalizeExtraction(extracted TaskExtraction) TaskExtraction {
	extracted.Title = strings.TrimSpace(extracted.Title)
	extracted.Priority = normalizePriority(extracted.Priority)
	extracted.Tags = normalizeTags(extracted.Tags)
	extracted.DueAt = normalizeOptionalString(extracted.DueAt)
	extracted.RemindAt = normalizeOptionalString(extracted.RemindAt)
	return extracted
}

func normalizePriority(raw string) string {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "low", "medium", "high":
		return strings.ToLower(strings.TrimSpace(raw))
	default:
		return "medium"
	}
}

func normalizeTags(tags []string) []string {
	if len(tags) == 0 {
		return []string{}
	}

	seen := make(map[string]struct{}, len(tags))
	normalized := make([]string, 0, len(tags))
	for _, tag := range tags {
		cleaned := strings.ToLower(strings.TrimSpace(tag))
		if cleaned == "" {
			continue
		}
		if _, exists := seen[cleaned]; exists {
			continue
		}
		seen[cleaned] = struct{}{}
		normalized = append(normalized, cleaned)
	}

	if len(normalized) == 0 {
		return []string{}
	}
	return normalized
}

func normalizeOptionalString(value *string) *string {
	if value == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}

func readOutputText(outputs []responsesOutputItem) string {
	for _, output := range outputs {
		if output.Type != "message" {
			continue
		}

		for _, content := range output.Content {
			if (content.Type == "output_text" || content.Type == "text") && strings.TrimSpace(content.Text) != "" {
				return content.Text
			}
		}
	}
	return ""
}

func truncateForError(value string, maxLen int) string {
	trimmed := strings.TrimSpace(value)
	if len(trimmed) <= maxLen {
		return trimmed
	}
	return trimmed[:maxLen] + "..."
}

func parseOpenAIErrorMessage(raw []byte) string {
	var payload struct {
		Error struct {
			Message string `json:"message"`
			Type    string `json:"type"`
			Code    string `json:"code"`
		} `json:"error"`
	}

	if err := json.Unmarshal(raw, &payload); err != nil {
		return ""
	}

	message := strings.TrimSpace(payload.Error.Message)
	if message == "" {
		return ""
	}

	if payload.Error.Code != "" {
		return fmt.Sprintf("%s (code: %s)", message, payload.Error.Code)
	}

	if payload.Error.Type != "" {
		return fmt.Sprintf("%s (type: %s)", message, payload.Error.Type)
	}

	return message
}

func taskExtractionJSONSchema() map[string]any {
	nullableString := map[string]any{
		"anyOf": []any{
			map[string]any{"type": "string"},
			map[string]any{"type": "null"},
		},
	}

	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"title": map[string]any{
				"type":      "string",
				"minLength": 1,
			},
			"due_at": nullableString,
			"priority": map[string]any{
				"type": "string",
				"enum": []string{"low", "medium", "high"},
			},
			"tags": map[string]any{
				"type": "array",
				"items": map[string]any{
					"type": "string",
				},
			},
			"remind_at": nullableString,
		},
		"required":             []string{"title", "due_at", "priority", "tags", "remind_at"},
		"additionalProperties": false,
	}
}

const taskExtractionSystemPrompt = `You extract todo task fields from natural language.

Rules:
- Return strict JSON matching the schema.
- title must be concise and action-oriented.
- due_at and remind_at must be RFC3339 datetime strings with timezone offset when present.
- If no due date is provided, return due_at as null.
- If no reminder is requested, return remind_at as null.
- If a reminder offset is requested and due_at is known, calculate remind_at.
- priority must be low, medium, or high.
- tags should be short lowercase labels and can be empty.`

type responsesRequest struct {
	Model string                  `json:"model"`
	Input []responsesInputMessage `json:"input"`
	Text  responsesTextConfig     `json:"text"`
}

type responsesInputMessage struct {
	Role    string                  `json:"role"`
	Content []responsesInputContent `json:"content"`
}

type responsesInputContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type responsesTextConfig struct {
	Format responsesTextFormat `json:"format"`
}

type responsesTextFormat struct {
	Type   string         `json:"type"`
	Name   string         `json:"name,omitempty"`
	Schema map[string]any `json:"schema,omitempty"`
	Strict bool           `json:"strict,omitempty"`
}

type responsesResponse struct {
	OutputText string                `json:"output_text"`
	Output     []responsesOutputItem `json:"output"`
	Error      *responsesAPIError    `json:"error,omitempty"`
}

type responsesOutputItem struct {
	Type    string                   `json:"type"`
	Content []responsesOutputContent `json:"content,omitempty"`
}

type responsesOutputContent struct {
	Type string `json:"type"`
	Text string `json:"text,omitempty"`
}

type responsesAPIError struct {
	Message string `json:"message"`
}
