package http_handler

import (
	"fmt"
	"gigaAPI/internal/logger"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

func newModuleResponse() *moduleResponse {
	return &moduleResponse{
		lg: logger.NewLogger(),
	}
}

type moduleResponse struct {
	lg *logger.Logger
}

func hashPass(p string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("hash pass error: ", err)
		return p
	}

	return string(hash)
}

func hashCompare(hash, pass string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pass))
	return err == nil
}

type SuccessResponse struct {
	Operation string      `json:"operation"`
	Result    string      `json:"result"`
	Data      interface{} `json:"data,omitempty"`
}

type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Error   string `json:"error"`
}

func (m *moduleResponse) SuccessCreated(c echo.Context, data interface{}) error {
	return c.JSON(http.StatusCreated, SuccessResponse{
		Operation: "create",
		Result:    "success",
		Data:      data,
	})
}

func (m *moduleResponse) SuccessLogout(c echo.Context, data interface{}) error {
	return c.JSON(http.StatusCreated, SuccessResponse{
		Operation: "logout",
		Result:    "success",
		Data:      data,
	})
}

func (m *moduleResponse) Error400(c echo.Context, msg string, e error) error {
	m.lg.LogServerError(e, "", c.Request().URL.String())
	return c.JSON(http.StatusBadRequest, ErrorResponse{
		Status:  http.StatusBadRequest,
		Message: msg,
		Error:   e.Error(),
	})
}

func (m *moduleResponse) Error500(c echo.Context, msg string, e error) error {
	m.lg.LogServerError(e, "", c.Request().URL.String())
	return c.JSON(http.StatusInternalServerError, ErrorResponse{
		Status:  http.StatusInternalServerError,
		Message: msg,
		Error:   e.Error(),
	})
}
