package http_server

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func (s *ServerHTTPS) rootHandle(c echo.Context) error {
	return c.Render(http.StatusOK, "root.html", nil)
}

func (s *ServerHTTPS) rootPath() {
	s.echo.GET("/worker", s.rootHandle)
}
