package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mrz1836/go-sanitize"
	"github.com/nyhryan/go-tiny-url/db"
	"github.com/nyhryan/go-tiny-url/utility"
	"github.com/nyhryan/go-tiny-url/views"
	"gorm.io/gorm"
)

type Handler struct {
	Db *gorm.DB
}

// GET / : Return the index page(views/index.templ)
func (h Handler) GetIndex(ctx echo.Context) error {
	return utility.Render(ctx, http.StatusOK, views.Index())
}

// GET /api : Query all records in db and render them into views/records.templ
// and return the rendered html, which is the list of records
func (h Handler) GetAPI(ctx echo.Context) error {
	var records []db.URLRecords
	h.Db.Find(&records)

	return utility.Render(ctx, http.StatusOK, views.Records(records))
}

// POST /api : Create a new record in db and render all records into views/records.templ
func (h Handler) PostAPI(ctx echo.Context) error {
	longURL := ctx.FormValue("longURL")
	longURL = sanitize.XSS(longURL)
	err := ctx.Validate(longURL)
	if err != nil {
		return utility.Render(ctx, http.StatusBadRequest, views.InvalidURL())
	}

	var tempRecords []db.URLRecords
	result := h.Db.Where("long_url = ?", longURL).Find(&tempRecords)

	// if longURL already exists in the database
	if result.RowsAffected != 0 {
		return utility.Render(ctx, http.StatusNotAcceptable, views.AlreadyExists())
	}

	tinyURL, err := utility.GenerateTinyURL(longURL, h.Db)
	if err != nil {
		ctx.Logger().Fatal(err)
	}

	record := &db.URLRecords{
		LongURL:    longURL,
		TinyURL:    tinyURL,
		ClickCount: 0,
	}

	h.Db.Create(record)

	var records []db.URLRecords
	h.Db.Find(&records)

	return utility.Render(ctx, http.StatusOK, views.Records(records))
}

// DELETE /api/:id : Delete a record in db
func (h Handler) DeleteAPI(ctx echo.Context) error {
	id := ctx.Param("id")
	h.Db.Delete(&db.URLRecords{}, id)

	return ctx.HTML(http.StatusOK, "<div hx-swap-oob='true' id='error'></div>")
}

// GET /:tinyURL : Redirect to the long URL
func (h Handler) RedirectTinyURL(ctx echo.Context) error {
	tinyURL := ctx.Param("tinyURL")
	if tinyURL == "favicon.ico" {
		return ctx.NoContent(http.StatusNoContent)
	}

	var record db.URLRecords
	result := h.Db.Where("tiny_url = ?", tinyURL).First(&record)
	if result.Error != nil {
		ctx.Logger().Fatal(result.Error)
	}
	record.ClickCount++

	h.Db.Model(&record).Update("click_count", record.ClickCount)

	return ctx.Redirect(http.StatusTemporaryRedirect, record.LongURL)
}
