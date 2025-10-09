package main

import (
	"context"
	"embed"
	"fmt"
	"io"
	"io/fs"
	"log"
	"mime"
	"os"
	"path"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Embed the built frontend so the binary can serve static files after deployment.
//go:embed client/dist/**
var embeddedFiles embed.FS


type Todo struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Completed bool               `json:"completed"`
	Body      string             `json:"body"`
}

var collection *mongo.Collection

func main() {
	fmt.Println("🚀 Starting Go Fiber server...")

	// Load .env only locally (ignore errors in production)
	_ = godotenv.Load(".env")

	MONGODB_URI := os.Getenv("MONGODB_URI")
	if MONGODB_URI == "" {
		log.Fatal("MONGODB_URI not set")
	}

	clientOptions := options.Client().ApplyURI(MONGODB_URI)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.Background())

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("✅ Connected to MongoDB Atlas")

	collection = client.Database("golang_db").Collection("todos")

	app := fiber.New()

	// ---- Embedded static files (client/dist) ----
	// Create a sub FS rooted at client/dist
	distFS, err := fs.Sub(embeddedFiles, "client/dist")
	if err != nil {
		log.Printf("⚠️  Could not create embedded fs: %v (static files may not be available at compile-time)", err)

		// Fallback: serve files from disk at runtime if `client/dist` exists.
		// This helps deployments that build the frontend during the deploy step
		// (producing `client/dist` on the host) but did not embed files at compile time.
		app.Static("/", "./client/dist")

		// SPA fallback: send index.html for client-side routes
		app.Get("/*", func(c *fiber.Ctx) error {
			if strings.HasPrefix(c.Path(), "/api/") {
				return c.Next()
			}
			reqPath := c.Path()
			if reqPath == "/" || reqPath == "" {
				reqPath = "/index.html"
			}
			fp := "./client/dist" + reqPath
			if _, statErr := os.Stat(fp); statErr != nil {
				return c.Status(404).SendString("Not found")
			}
			return c.SendFile(fp)
		})
	} else {
		// Serve embedded assets and provide SPA fallback
		app.Get("/*", func(c *fiber.Ctx) error {
			// Let API routes pass through
			if strings.HasPrefix(c.Path(), "/api/") {
				return c.Next()
			}

			reqPath := c.Path()
			if reqPath == "/" || reqPath == "" {
				reqPath = "/index.html"
			}
			fp := strings.TrimPrefix(reqPath, "/")

			// Try to open the requested file from embedded FS
			f, ferr := distFS.Open(fp)
			if ferr != nil {
				// Fallback to index.html for client-side routing
				idx, ierr := distFS.Open("index.html")
				if ierr != nil {
					return c.Status(404).SendString("Not found")
				}
				defer idx.Close()
				c.Set("Content-Type", "text/html; charset=utf-8")
				return c.SendStream(idx)
			}
			defer f.Close()

			// Set content-type based on extension when available
			if ext := path.Ext(fp); ext != "" {
				if m := mime.TypeByExtension(ext); m != "" {
					c.Set("Content-Type", m)
				}
			}

			// Stream the file
			return c.SendStream(io.NopCloser(f))
		})
	}

	// Health check for Render
	app.Get("/healthz", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	// Enable CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:5174, https://todoapp-frontend.onrender.com",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	// Routes
	app.Get("/api/todos", getTodos)
	app.Post("/api/todos", createTodo)
	app.Patch("/api/todos/:id", updateTodo)
	app.Delete("/api/todos/:id", deleteTodo)

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}
	log.Printf("✅ Server running on port %s", port)
	log.Fatal(app.Listen("0.0.0.0:" + port))
}

func getTodos(c *fiber.Ctx) error {
	var todos []Todo
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		return err
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var todo Todo
		if err := cursor.Decode(&todo); err != nil {
			return err
		}
		todos = append(todos, todo)
	}
	return c.JSON(todos)
}

func createTodo(c *fiber.Ctx) error {
	todo := new(Todo)
	if err := c.BodyParser(todo); err != nil {
		return err
	}
	if todo.Body == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Todo body cannot be empty"})
	}
	insertResult, err := collection.InsertOne(context.Background(), todo)
	if err != nil {
		return err
	}
	todo.ID = insertResult.InsertedID.(primitive.ObjectID)
	return c.Status(201).JSON(todo)
}

func updateTodo(c *fiber.Ctx) error {
	id := c.Params("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid todo ID"})
	}
	filter := bson.M{"_id": objectID}
	update := bson.M{"$set": bson.M{"completed": true}}
	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}
	return c.Status(200).JSON(fiber.Map{"success": true})
}

func deleteTodo(c *fiber.Ctx) error {
	id := c.Params("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid todo ID"})
	}
	filter := bson.M{"_id": objectID}
	_, err = collection.DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}
	return c.Status(200).JSON(fiber.Map{"success": true})
}
