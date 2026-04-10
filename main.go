package main

import (
	"context"
	"embed"
	"errors"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/krish-anu/ToDoAppBackend/internal/ai"
	"github.com/krish-anu/ToDoAppBackend/internal/config"
	"github.com/krish-anu/ToDoAppBackend/internal/database"
	"github.com/krish-anu/ToDoAppBackend/internal/handlers"
	"github.com/krish-anu/ToDoAppBackend/internal/routes"
	"github.com/krish-anu/ToDoAppBackend/internal/static"
	"github.com/krish-anu/ToDoAppBackend/internal/workflow"
)

// Embed the built frontend so the binary can serve static files after deployment.
//
//go:embed client/dist/**
var embeddedFiles embed.FS

func main() {
	cfg, err := config.Load()
	if err != nil {
		if errors.Is(err, config.ErrMissingMongoURI) {
			log.Fatal(err)
		}
		log.Fatalf("failed to load configuration: %v", err)
	}

	ctx := context.Background()
	mongoClient, todoCollection, err := database.Connect(ctx, cfg.MongoURI, cfg.DatabaseName, cfg.CollectionName)
	if err != nil {
		log.Fatalf("failed to connect to mongo: %v", err)
	}
	defer mongoClient.Disconnect(ctx)

	databaseRef := mongoClient.Database(cfg.DatabaseName)
	reminderCollection := databaseRef.Collection(cfg.ReminderCollectionName)
	workflowRunsCollection := databaseRef.Collection(cfg.WorkflowRunsCollectionName)

	app := fiber.New()

	// Health check for Render
	app.Get("/healthz", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins: cfg.AllowOrigins,
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	if cfg.OpenAIAPIKey == "" {
		log.Println("OPENAI_API_KEY not set. /api/workflows/task-from-text will return 503 until configured.")
	}

	todoHandler := handlers.NewTodoHandler(todoCollection)
	taskExtractor := ai.NewOpenAITaskExtractor(cfg.OpenAIAPIKey, cfg.OpenAIModel, cfg.OpenAIBaseURL)
	taskFromTextWorkflow := workflow.NewTaskFromTextWorkflow(
		todoCollection,
		reminderCollection,
		workflowRunsCollection,
		taskExtractor,
		cfg.WorkflowDefaultUserID,
		cfg.WorkflowDefaultTimezone,
	)
	workflowHandler := handlers.NewWorkflowHandler(taskFromTextWorkflow)

	routes.RegisterTodoRoutes(app, todoHandler)
	routes.RegisterWorkflowRoutes(app, workflowHandler)
	static.Register(app, embeddedFiles)

	log.Printf("server running on port %s", cfg.Port)
	log.Fatal(app.Listen("0.0.0.0:" + cfg.Port))
}
