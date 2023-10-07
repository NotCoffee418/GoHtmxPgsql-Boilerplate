package validation

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type FormResponse struct {
	IsSuccessful bool
	Message      string
}

// RenderFormResponse provides renders a validation response, intended to be used with HTMX.
// Empty slice and nil are valid inputs, it will simply reserve the space for the response.
func RenderFormResponse(c *gin.Context, validationData []FormResponse) {
	c.HTML(http.StatusOK, "form_response.html", gin.H{
		"ValidationMessages": validationData,
	})
}
