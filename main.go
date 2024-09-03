package main

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/nyhryan/go-tiny-url/db"
	"github.com/nyhryan/go-tiny-url/handler"
	"github.com/nyhryan/go-tiny-url/utility"
)

func main() {
	// Create a new echo instance
	e := echo.New()
	e.Validator = &utility.URLValidator{Validator: validator.New(validator.WithRequiredStructEnabled())}

	// Create a new handler instance
	h := handler.Handler{}
	h.Db = db.New()

	// Routes

	// GET / - Get the index page
	e.GET("/", h.GetIndex)

	// GET /api - Get all records
	e.GET("/api", h.GetAPI)

	// POST /api - Create a new record
	e.POST("/api", h.PostAPI)

	// GET /:tinyURL - Redirect to the long URL
	e.GET("/:tinyURL", h.RedirectTinyURL)

	// DELETE /api/:id - Delete a record
	e.DELETE("/api/:id", h.DeleteAPI)

	// Start the server
	e.Logger.Fatal(e.Start("localhost:5000"))
}
