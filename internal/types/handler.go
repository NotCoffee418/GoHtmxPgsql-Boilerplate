package types

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type HandlerRegistrar interface {
	Handler(engine *gin.Engine)
	Initialize(initContext *HandlerInitContext)
}

// HandlerInitContext is passed to the Initialize method of a HandlerRegistrar
type HandlerInitContext struct {
	DB *sqlx.DB
}
