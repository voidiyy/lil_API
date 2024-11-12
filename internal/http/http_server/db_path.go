package http_server

import (
	"gigaAPI/internal/http/http_server/http_security"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (s *ServerHTTPS) dbPath() {
	//root
	s.echo.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "greet")
	})

	//registration
	s.echo.GET("/worker/register", func(c echo.Context) error {
		return c.Render(200, "register.html", nil)
	})
	s.echo.POST("/worker/register", s.handler.RegisterWorker)

	//login
	s.echo.GET("/worker/login", func(c echo.Context) error {
		return c.Render(200, "login.html", nil)
	}, http_security.SSetCSRF())
	s.echo.POST("/worker/login", s.handler.LoginWorker)

	//logout
	s.echo.GET("/worker/logout", s.handler.LogoutWorker)

	//profile
	s.echo.GET("/worker/profile", s.handler.ProfileWorker)
}
