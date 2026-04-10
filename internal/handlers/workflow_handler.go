package handlers

import (
	"errors"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/krish-anu/ToDoAppBackend/internal/ai"
	"github.com/krish-anu/ToDoAppBackend/internal/workflow"
)

type WorkflowHandler struct {
	taskFromTextWorkflow *workflow.TaskFromTextWorkflow
}

func NewWorkflowHandler(taskFromTextWorkflow *workflow.TaskFromTextWorkflow) *WorkflowHandler {
	return &WorkflowHandler{taskFromTextWorkflow: taskFromTextWorkflow}
}

type createTaskFromTextRequest struct {
	UserID   string `json:"user_id"`
	Message  string `json:"message"`
	Timezone string `json:"timezone"`
}

func (h *WorkflowHandler) CreateTaskFromText(c *fiber.Ctx) error {
	var req createTaskFromTextRequest
	if err := c.BodyParser(&req); err != nil {
		return sendError(c, fiber.StatusBadRequest, "Invalid request payload")
	}

	if strings.TrimSpace(req.Message) == "" {
		return sendError(c, fiber.StatusBadRequest, "message is required")
	}

	result, err := h.taskFromTextWorkflow.Run(c.UserContext(), workflow.TaskFromTextInput{
		UserID:   req.UserID,
		Message:  req.Message,
		Timezone: req.Timezone,
	})
	if err != nil {
		log.Printf("task-from-text workflow failed: %v", err)

		errMessage := err.Error()
		if strings.Contains(errMessage, "api key") || strings.Contains(errMessage, "OPENAI") {
			return sendError(c, fiber.StatusServiceUnavailable, "Workflow is not configured. Set OPENAI_API_KEY.")
		}

		if isUpstreamOpenAIError(errMessage) || errors.Is(err, ai.ErrMalformedResponse) {
			return sendError(c, fiber.StatusBadGateway, errMessage)
		}

		switch {
		case errors.Is(err, workflow.ErrInvalidMessage), errors.Is(err, workflow.ErrInvalidTimezone):
			return sendError(c, fiber.StatusBadRequest, err.Error())
		case errors.Is(err, ai.ErrMissingAPIKey):
			return sendError(c, fiber.StatusServiceUnavailable, "Workflow is not configured. Set OPENAI_API_KEY.")
		default:
			return sendError(c, fiber.StatusInternalServerError, errMessage)
		}
	}

	statusCode := fiber.StatusCreated
	if result.Partial {
		statusCode = fiber.StatusAccepted
	}

	return c.Status(statusCode).JSON(fiber.Map{
		"success":  true,
		"message":  result.Message,
		"partial":  result.Partial,
		"todo":     result.Todo,
		"reminder": result.Reminder,
		"parsed":   result.Parsed,
	})
}

func isUpstreamOpenAIError(message string) bool {
	value := strings.ToLower(strings.TrimSpace(message))
	return strings.Contains(value, "openai returned status") ||
		strings.Contains(value, "openai error") ||
		strings.Contains(value, "failed to call openai") ||
		strings.Contains(value, "failed to decode openai response") ||
		strings.Contains(value, "failed to read openai response") ||
		strings.Contains(value, "response missing structured output")
}
