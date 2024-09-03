package utility

import (
	"crypto/md5"
	"fmt"
	"io"
	"time"

	"github.com/a-h/templ"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/nyhryan/go-tiny-url/db"
	"gorm.io/gorm"
)

type URLValidator struct {
	Validator *validator.Validate
}

func (uv *URLValidator) Validate(i interface{}) error {
	err := uv.Validator.Var(i, "required,http_url")
	if err != nil {
		return err
	}
	return nil
}

// Render templ files into HTML
func Render(ctx echo.Context, statusCode int, t templ.Component) error {
	buf := templ.GetBuffer()
	defer templ.ReleaseBuffer(buf)

	err := t.Render(ctx.Request().Context(), buf)
	if err != nil {
		return err
	}

	return ctx.HTML(statusCode, buf.String())
}

// Generates tiny url based on MD5 hash and uses 7 characters of the hash
func GenerateTinyURL(longURL string, database *gorm.DB) (string, error) {
	hash := md5.New()
	io.WriteString(hash, longURL)
	io.WriteString(hash, fmt.Sprintf("%d", time.Now().Unix()))
	hashed := fmt.Sprintf("%x", hash.Sum(nil))

	var tinyURL string
	for i := 0; i < len(hashed); i++ {
		tinyURL = hashed[i : i+7]
		var tempRecords []db.URLRecords
		result := database.Where("tiny_url = ?", tinyURL).Find(&tempRecords)

		// if generated tinyURL is not in the database
		if result.RowsAffected == 0 {
			break
		} else {
			// if the last 7 characters of hash has been already used, throw an error
			if i == len(hashed)-7 {
				return "", fmt.Errorf("cannot generate a unique tinyURL")
			} else {
				continue
			}
		}
	}
	return tinyURL, nil
}
