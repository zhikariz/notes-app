package helper

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Response struct {
	Meta Meta        `json:"meta"`
	Data interface{} `json:"data"`
}

type Meta struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Status  string `json:"status"`
}

func APIResponse(message string, code int, status string, data interface{}) Response {
	meta := Meta{
		Message: message,
		Code:    code,
		Status:  status,
	}
	jsonResponse := Response{
		Meta: meta,
		Data: data,
	}
	return jsonResponse
}

func FormatValidationError(err error) []string {
	var errors []string

	for _, e := range err.(validator.ValidationErrors) {
		errors = append(errors, e.Error())
	}
	return errors

}

func ErrorValidation(c *gin.Context, err error, message string) {
	errors := FormatValidationError(err)
	errorMessage := gin.H{"errors": errors}

	response := APIResponse(message, http.StatusBadRequest, "error", errorMessage)
	c.JSON(http.StatusBadRequest, response)
}

func ErrorHandling(c *gin.Context, err error, message string) {
	var errors []string
	errors = append(errors, err.Error())
	errorMessage := gin.H{"errors": errors}
	response := APIResponse(message, http.StatusUnprocessableEntity, "error", errorMessage)
	c.JSON(http.StatusUnprocessableEntity, response)
}

func SuccessHandling(c *gin.Context, formatter interface{}, message string) {
	response := APIResponse(message, http.StatusOK, "success", formatter)

	c.JSON(http.StatusOK, response)
}

func AuthorizationHandling(c *gin.Context) {
	var errors []string
	errors = append(errors, "You're Not Authorized to do this")
	errorMessage := gin.H{"errors": errors}

	response := APIResponse("Unauthorized", http.StatusUnauthorized, "error", errorMessage)
	c.AbortWithStatusJSON(http.StatusUnauthorized, response)

}
