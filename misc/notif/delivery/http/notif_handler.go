package http

import (
	"context"
	"net/http"

	"github.com/labstack/echo"
	"github.com/misc/notif"
	"github.com/models"
	"github.com/sirupsen/logrus"
)

type ResponseError struct {
	Message string `json:"message"`
}

type notifHandler struct {
	notifUsecase notif.Usecase
}

func NewNotifHandler(e *echo.Echo, us notif.Usecase) {
	handler := &notifHandler{
		notifUsecase: us,
	}
	e.GET("misc/notif", handler.Get)
}

func (a *notifHandler) Get(c echo.Context) error {
	c.Request().Header.Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	c.Response().Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	token := c.Request().Header.Get("Authorization")

	if token == "" {
		return c.JSON(http.StatusUnauthorized, models.ErrUnAuthorize)
	}

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	res, err := a.notifUsecase.GetByMerchantID(ctx, token)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, res)
}

func getStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}
	logrus.Error(err)
	switch err {
	case models.ErrInternalServerError:
		return http.StatusInternalServerError
	case models.ErrNotFound:
		return http.StatusNotFound
	case models.ErrUnAuthorize:
		return http.StatusUnauthorized
	case models.ErrConflict:
		return http.StatusBadRequest
	case models.ErrBadParamInput:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}