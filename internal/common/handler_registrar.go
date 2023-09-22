package common

import (
	"github.com/gin-gonic/gin"
)

type HandlerRegistrar interface {
	Handler(engine *gin.Engine)
}
