package static

import (
	"embed"
	"io"
	"io/fs"
	"log"
	"mime"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func Register(app *fiber.App, embeddedFiles embed.FS) {
	distFS, err := fs.Sub(embeddedFiles, "web/dist")
	if err != nil {
		log.Printf("could not create embedded fs: %v (falling back to runtime static files)", err)
		registerDiskFallback(app)
		return
	}

	app.Get("/*", func(c *fiber.Ctx) error {
		if strings.HasPrefix(c.Path(), "/api/") {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Route not found"})
		}

		reqPath := c.Path()
		if reqPath == "/" || reqPath == "" {
			reqPath = "/index.html"
		}

		fp := strings.TrimPrefix(reqPath, "/")
		file, openErr := distFS.Open(fp)
		if openErr != nil {
			indexFile, indexErr := distFS.Open("index.html")
			if indexErr != nil {
				return c.Status(fiber.StatusNotFound).SendString("Not found")
			}
			defer indexFile.Close()
			c.Set("Content-Type", "text/html; charset=utf-8")
			return c.SendStream(indexFile)
		}
		defer file.Close()

		if ext := path.Ext(fp); ext != "" {
			if contentType := mime.TypeByExtension(ext); contentType != "" {
				c.Set("Content-Type", contentType)
			}
		}

		return c.SendStream(io.NopCloser(file))
	})
}

func registerDiskFallback(app *fiber.App) {
	distDir := resolveDistDir()
	app.Static("/", distDir)

	app.Get("/*", func(c *fiber.Ctx) error {
		if strings.HasPrefix(c.Path(), "/api/") {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Route not found"})
		}

		reqPath := c.Path()
		if reqPath == "/" || reqPath == "" {
			reqPath = "/index.html"
		}

		filePath := filepath.Join(distDir, strings.TrimPrefix(reqPath, "/"))
		if _, statErr := os.Stat(filePath); statErr != nil {
			return c.Status(fiber.StatusNotFound).SendString("Not found")
		}

		return c.SendFile(filePath)
	})
}

func resolveDistDir() string {
	candidates := []string{
		"./web/dist",
		"./backend/web/dist",
		"../frontend/dist",
	}

	for _, candidate := range candidates {
		if _, err := os.Stat(filepath.Join(candidate, "index.html")); err == nil {
			return candidate
		}
	}

	return "./web/dist"
}
