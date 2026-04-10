package main

import (
	"context"
	"embed"
	"errors"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/krish-anu/ToDoAppBackend/internal/config"
	"github.com/krish-anu/ToDoAppBackend/internal/database"
	"github.com/krish-anu/ToDoAppBackend/internal/handlers"
	"github.com/krish-anu/ToDoAppBackend/internal/routes"
	"github.com/krish-anu/ToDoAppBackend/internal/static"
)

// Embed the built frontend so the binary can serve static files after deployment.
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
	mongoClient, collection, err := database.Connect(ctx, cfg.MongoURI, cfg.DatabaseName, cfg.CollectionName)
	if err != nil {
		log.Fatalf("failed to connect to mongo: %v", err)
	}
	defer mongoClient.Disconnect(ctx)

	app := fiber.New()

	// Health check for Render
	app.Get("/healthz", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins: cfg.AllowOrigins,
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	todoHandler := handlers.NewTodoHandler(collection)
	routes.RegisterTodoRoutes(app, todoHandler)
	static.Register(app, embeddedFiles)

	log.Printf("server running on port %s", cfg.Port)
	log.Fatal(app.Listen("0.0.0.0:" + cfg.Port))
}
