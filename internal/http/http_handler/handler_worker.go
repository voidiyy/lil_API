package http_handler

import (
	"fmt"
	"gigaAPI/internal/db_boss"
	"gigaAPI/internal/http/http_server/http_security"
	"gigaAPI/internal/logger"
	types "gigaAPI/internal/type"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
)

const path = "http/http_handler/handler_worker."

func tracer(s string) string {
	return path + s
}

var _ HandlerWorker = &handlerW{}

type handlerW struct {
	db        *db_boss.HandlerDB
	lg        *logger.Logger
	resp      *moduleResponse
	jwtSecret []byte
}

func NewHandlerWorker(db *db_boss.HandlerDB, jwtSecret []byte) HandlerWorker {
	return &handlerW{
		db:        db,
		lg:        logger.NewLogger(),
		resp:      newModuleResponse(),
		jwtSecret: jwtSecret,
	}
}

type HandlerWorker interface {
	RegisterWorker(c echo.Context) error
	LoginWorker(c echo.Context) error
	ProfileWorker(c echo.Context) error
	LogoutWorker(c echo.Context) error
}

func (h *handlerW) ProfileWorker(c echo.Context) error {
	claims, err := http_security.GetJWTCookie(c, h.jwtSecret)
	if err != nil || claims == nil {
		return c.Render(http.StatusOK, "profile.html", map[string]interface{}{
			"Error": "You are not authorized. Please log in to view your profile.",
		})
	}

	ctx := c.Request().Context()
	worker, er := h.db.Worker.GetW(ctx, claims.ID)
	if er != nil {
		return h.resp.Error500(c, "error get worker", er)
	}

	return c.Render(http.StatusOK, "profile.html", map[string]interface{}{
		"Name":      worker.Name,
		"Email":     worker.Email,
		"Role":      worker.Role,
		"CreatedAt": worker.CreatedAt.Format("2006-01-02"),
	})
}

func (h *handlerW) LogoutWorker(c echo.Context) error {
	http_security.DeleteJWTCookie(c)
	return h.resp.SuccessLogout(c, nil)
}

func (h *handlerW) LoginWorker(c echo.Context) error {
	var (
		err  error
		form LoginWorkerRequest
	)

	err = c.Bind(&form)
	if err != nil {
		return h.resp.Error400(c, "read form error", err)
	}

	err = c.Validate(&form)
	if err != nil {
		return h.resp.Error400(c, "invalid login data", err)
	}

	ctx := c.Request().Context()

	worker, er := h.db.Worker.GetWByEmail(ctx, form.Email)
	if er != nil {
		return h.resp.Error500(c, er.Error(), er)
	}

	if !hashCompare(worker.Password, form.Password) {
		return h.resp.Error400(c, "invalid password", fmt.Errorf("password %s invalid for this accout", form.Password))
	}

	err = http_security.SaveJWTCookie(c, worker.ID, h.jwtSecret)
	if err != nil {
		return h.resp.Error400(c, "error generate JWT", err)
	}

	return c.JSON(200, nil)
}

func (h *handlerW) RegisterWorker(c echo.Context) error {
	var (
		err  error
		form RegisterWorkerRequest
	)

	err = c.Bind(&form)
	if err != nil {
		return h.resp.Error400(c, "read form error", err)
	}

	err = c.Validate(&form)
	if err != nil {
		return h.resp.Error400(c, "validate form error", err)
	}

	hash := hashPass(form.Password)

	w := &types.Worker{
		ID:       uuid.New(),
		Name:     form.Name,
		Password: hash,
		Email:    form.Email,
		Role:     form.Role,
		IsActive: true,
	}
	ctx := c.Request().Context()
	err = h.db.Worker.CreateW(ctx, w)
	if err != nil {
		return h.resp.Error500(c, err.Error(), err)
	}

	return h.resp.SuccessCreated(c, w.Name)
}
