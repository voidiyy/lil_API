package http_server

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"html/template"
	"io"
)

type HTML struct {
	templates *template.Template
}

func NewHTML(templateDir string) (*HTML, error) {
	tmpl, err := template.ParseGlob(templateDir + "/*/*.html")
	if err != nil {
		return nil, fmt.Errorf("failed to parse templates: %w", err)
	}
	return &HTML{
		templates: tmpl,
	}, nil
}

func (h *HTML) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return h.templates.ExecuteTemplate(w, name, data)
}
