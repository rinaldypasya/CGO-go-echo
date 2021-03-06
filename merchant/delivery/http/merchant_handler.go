package http

import (
	"context"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	validator "gopkg.in/go-playground/validator.v9"

	"github.com/merchant"
	"github.com/models"
)

// ResponseError represent the reseponse error struct
type ResponseError struct {
	Message string `json:"message"`
}

// merchantHandler  represent the httphandler for merchant
type merchantHandler struct {
	MerchantUsecase merchant.Usecase
}

// NewmerchantHandler will initialize the merchants/ resources endpoint
func NewmerchantHandler(e *echo.Echo, us merchant.Usecase) {
	handler := &merchantHandler{
		MerchantUsecase: us,
	}
	e.POST("/merchants", handler.CreateMerchant)
	e.PUT("/merchants/:id", handler.UpdateMerchant)
	//e.GET("/merchants/:id", handler.GetByID)
	//e.DELETE("/merchants/:id", handler.Delete)
}

func isRequestValid(m *models.NewCommandMerchant) (bool, error) {
	validate := validator.New()
	err := validate.Struct(m)
	if err != nil {
		return false, err
	}
	return true, nil
}

// Store will store the merchant by given request body
func (a *merchantHandler) CreateMerchant(c echo.Context) error {
	//var merchantCommand models.NewCommandMerchant
	//err := c.Bind(&merchantCommand)
	//if err != nil {
	//	return c.JSON(http.StatusUnprocessableEntity, err.Error())
	//}
	balance, _ := strconv.ParseFloat(c.FormValue("balance"),64)
	merchantCommand := models.NewCommandMerchant{
		Id:               c.FormValue("id"),
		MerchantName:     c.FormValue("merchant_name"),
		MerchantDesc:     c.FormValue("merchant_desc"),
		MerchantEmail:    c.FormValue("merchant_email"),
		MerchantPassword: c.FormValue("password"),
		Balance:          balance,
	}
	if ok, err := isRequestValid(&merchantCommand); !ok {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	error := a.MerchantUsecase.Create(ctx, &merchantCommand,"admin")

	if error != nil {
		return c.JSON(getStatusCode(error), ResponseError{Message: error.Error()})
	}
	return c.JSON(http.StatusCreated, merchantCommand)
}

func (a *merchantHandler) UpdateMerchant(c echo.Context) error {
	//var merchantCommand models.NewCommandMerchant
	//err := c.Bind(&merchantCommand)
	//if err != nil {
	//	return c.JSON(http.StatusUnprocessableEntity, err.Error())
	//}
	balance, _ := strconv.ParseFloat(c.FormValue("balance"),64)
	merchantCommand := models.NewCommandMerchant{
		Id:               c.FormValue("id"),
		MerchantName:     c.FormValue("merchant_name"),
		MerchantDesc:     c.FormValue("merchant_desc"),
		MerchantEmail:    c.FormValue("merchant_email"),
		MerchantPassword: c.FormValue("password"),
		Balance:          balance,
	}
	if ok, err := isRequestValid(&merchantCommand); !ok {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	err := a.MerchantUsecase.Update(ctx, &merchantCommand,"admin")

	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusCreated, merchantCommand)
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
		return http.StatusConflict
	case models.ErrBadParamInput:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
