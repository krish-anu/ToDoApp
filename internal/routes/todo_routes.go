package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/krish-anu/ToDoAppBackend/internal/handlers"
)

func RegisterTodoRoutes(app *fiber.App, todoHandler *handlers.TodoHandler) {
	app.Get("/api/todos", todoHandler.GetTodos)
	app.Post("/api/todos", todoHandler.CreateTodo)
	app.Patch("/api/todos/:id", todoHandler.UpdateTodo)
	app.Delete("/api/todos/:id", todoHandler.DeleteTodo)
}
