package handlers

import (
	"context"
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/krish-anu/ToDoAppBackend/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TodoHandler struct {
	collection *mongo.Collection
}

func NewTodoHandler(collection *mongo.Collection) *TodoHandler {
	return &TodoHandler{collection: collection}
}

func (h *TodoHandler) GetTodos(c *fiber.Ctx) error {
	ctx := context.Background()

	cursor, err := h.collection.Find(ctx, bson.M{}, options.Find().SetSort(bson.D{{Key: "_id", Value: -1}}))
	if err != nil {
		return sendError(c, fiber.StatusInternalServerError, "Failed to fetch todos")
	}
	defer cursor.Close(ctx)

	var todos []models.Todo
	if err := cursor.All(ctx, &todos); err != nil {
		return sendError(c, fiber.StatusInternalServerError, "Failed to decode todos")
	}

	if todos == nil {
		todos = []models.Todo{}
	}

	return c.JSON(todos)
}

func (h *TodoHandler) CreateTodo(c *fiber.Ctx) error {
	var req models.CreateTodoRequest
	if err := c.BodyParser(&req); err != nil {
		return sendError(c, fiber.StatusBadRequest, "Invalid request payload")
	}

	body := strings.TrimSpace(req.Body)
	if body == "" {
		return sendError(c, fiber.StatusBadRequest, "Todo body cannot be empty")
	}

	todo := models.Todo{
		Body:      body,
		Completed: req.Completed,
	}

	insertResult, err := h.collection.InsertOne(context.Background(), todo)
	if err != nil {
		return sendError(c, fiber.StatusInternalServerError, "Failed to create todo")
	}

	insertedID, ok := insertResult.InsertedID.(primitive.ObjectID)
	if !ok {
		return sendError(c, fiber.StatusInternalServerError, "Failed to parse inserted id")
	}

	todo.ID = insertedID
	return c.Status(fiber.StatusCreated).JSON(todo)
}

func (h *TodoHandler) UpdateTodo(c *fiber.Ctx) error {
	objectID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return sendError(c, fiber.StatusBadRequest, "Invalid todo ID")
	}

	var req models.UpdateTodoRequest
	if err := c.BodyParser(&req); err != nil {
		return sendError(c, fiber.StatusBadRequest, "Invalid request payload")
	}

	if req.Completed == nil {
		return sendError(c, fiber.StatusBadRequest, "completed is required")
	}

	filter := bson.M{"_id": objectID}
	update := bson.M{"$set": bson.M{"completed": *req.Completed}}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updated models.Todo
	err = h.collection.FindOneAndUpdate(context.Background(), filter, update, opts).Decode(&updated)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return sendError(c, fiber.StatusNotFound, "Todo not found")
	}
	if err != nil {
		return sendError(c, fiber.StatusInternalServerError, "Failed to update todo")
	}

	return c.JSON(updated)
}

func (h *TodoHandler) DeleteTodo(c *fiber.Ctx) error {
	objectID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return sendError(c, fiber.StatusBadRequest, "Invalid todo ID")
	}

	result, err := h.collection.DeleteOne(context.Background(), bson.M{"_id": objectID})
	if err != nil {
		return sendError(c, fiber.StatusInternalServerError, "Failed to delete todo")
	}
	if result.DeletedCount == 0 {
		return sendError(c, fiber.StatusNotFound, "Todo not found")
	}

	return c.JSON(fiber.Map{"success": true})
}

func sendError(c *fiber.Ctx, statusCode int, message string) error {
	return c.Status(statusCode).JSON(fiber.Map{"error": message})
}
