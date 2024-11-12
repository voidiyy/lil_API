package http_security

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
)

func SSetCSRF() echo.MiddlewareFunc {
	return middleware.CSRFWithConfig(middleware.CSRFConfig{
		TokenLookup:    "header:" + echo.HeaderXCSRFToken,
		CookieName:     "_csrf",
		CookiePath:     "/",
		CookieSecure:   true,
		CookieHTTPOnly: true,
		CookieSameSite: http.SameSiteStrictMode,
		ContextKey:     "_csrf",
	})
}

func CSRFTokenHeaderInjector() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			csrfCookie, err := c.Cookie("_csrf")
			if err == nil && csrfCookie != nil && csrfCookie.Value != "" {
				c.Request().Header.Set(echo.HeaderXCSRFToken, csrfCookie.Value)
			}
			return next(c)
		}
	}
}

func VerifyCSRF() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			srvToken := c.Get("_csrf")

			clientCookie, err := c.Cookie("_csrf")
			if err != nil || clientCookie.Value == "" {
				return c.JSON(403, map[string]string{
					"status": "unauthorized",
					"error":  "missing CSRF token in cookie",
				})
			}
			clientToken := clientCookie.Value

			fmt.Println("verify CSRF:")
			fmt.Println("srv token:", srvToken)
			fmt.Println("client token:", clientToken)
			if clientToken != srvToken {
				return c.JSON(404, map[string]string{
					"status": "unauthorized",
					"error":  "CSRF token mismatch",
				})
			}
			return next(c)
		}
	}
}
