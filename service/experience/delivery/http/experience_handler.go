package http

import (
	"github.com/labstack/echo"
	"github.com/models"
	"github.com/service/experience"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	validator "gopkg.in/go-playground/validator.v9"
	"net/http"
	"strconv"
)

// ResponseError represent the reseponse error struct
type ResponseError struct {
	Message string `json:"message"`
}

// experienceHandler  represent the httphandler for experience
type experienceHandler struct {
	experienceUsecase experience.Usecase
}

// NewexperienceHandler will initialize the experiences/ resources endpoint
func NewexperienceHandler(e *echo.Echo, us experience.Usecase) {
	handler := &experienceHandler{
		experienceUsecase: us,
	}
	//e.POST("/experiences", handler.Createexperience)
	//e.PUT("/experiences/:id", handler.Updateexperience)
	e.GET("service/experience/:id", handler.GetByID)
	e.GET("service/experience/search", handler.SearchExp)
	e.GET("service/experience/get-user-discover-preference", handler.GetUserDiscoverPreference)
	//e.DELETE("/experiences/:id", handler.Delete)
}

func isRequestValid(m *models.NewCommandMerchant) (bool, error) {
	validate := validator.New()
	err := validate.Struct(m)
	if err != nil {
		return false, err
	}
	return true, nil
}

// GetByID will get article by given id
func (a *experienceHandler) GetByID(c echo.Context) error {
	id := c.Param("id")
	//if err != nil {
	//	return c.JSON(http.StatusNotFound, models.ErrNotFound.Error())
	//}

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	art, err := a.experienceUsecase.GetByID(ctx, id)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, art)
}

func (a *experienceHandler) SearchExp(c echo.Context) error {
	harborID := c.QueryParam("harbor_id")
	cityID := c.QueryParam("city_id")

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	searchResult, err := a.experienceUsecase.SearchExp(ctx, harborID, cityID)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, searchResult)
}
func (a *experienceHandler) GetUserDiscoverPreference(c echo.Context) error {
	qpage := c.QueryParam("page")
	qsize := c.QueryParam("size")

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	if qpage != "" && qsize != ""{
		page , _:= strconv.Atoi(qpage)
		size , _:= strconv.Atoi(qsize)
		art, err := a.experienceUsecase.GetUserDiscoverPreference(ctx,&size,&page)
		if err != nil {
			return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
		}
		return c.JSON(http.StatusOK, art)
	}else {
		art, err := a.experienceUsecase.GetUserDiscoverPreference(ctx,nil,nil)
		if err != nil {
			return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
		}
		return c.JSON(http.StatusOK, art)
	}

	return nil
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
