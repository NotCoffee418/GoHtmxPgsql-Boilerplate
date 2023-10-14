package weberrors

import (
	"errors"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// InternalServerErrorResponse is a helper function to log an internal server error and return a 500 response.
// err is logged but not displayed to the user.
func InternalServerErrorResponse(c *gin.Context, err error) {
	log.Errorln("500: ", ExtractErrMsg(err))
	c.Error(gin.Error{
		Err:  errors.New("Internal Server Error"),
		Type: gin.ErrorTypePrivate,
	})
}

// ForbiddenResponse logs and sets a 403 Forbidden Error
func ForbiddenResponse(c *gin.Context) {
	log.Errorln("403: " + c.RemoteIP() + " - " + c.Request.URL.Path)
	c.Error(gin.Error{
		Err:  errors.New("Forbidden"),
		Type: gin.ErrorTypePrivate,
	})
}

// BadRequestResponse logs and sets a 400 Bad Request Error
// err is logged but not displayed to the user.
func BadRequestResponse(c *gin.Context, err error) {
	log.Errorln("400 on path " + c.Request.URL.Path + ": " + ExtractErrMsg(err))
	c.Error(gin.Error{
		Err:  errors.New("Bad Request"),
		Type: gin.ErrorTypePrivate,
	})
}

func ExtractErrMsg(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}
