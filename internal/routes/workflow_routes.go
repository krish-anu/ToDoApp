package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/krish-anu/ToDoAppBackend/internal/handlers"
)

func RegisterWorkflowRoutes(app *fiber.App, workflowHandler *handlers.WorkflowHandler) {
	app.Post("/api/workflows/task-from-text", workflowHandler.CreateTaskFromText)
}
