package types

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type HandlerRegistrar interface {
	Handler(engine *gin.Engine, db *sqlx.DB)
}
