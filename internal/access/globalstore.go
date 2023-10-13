package access

import (
	"html/template"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

// This should probably be refactored. It will do for the demo.

// Tmpl is the global template variable
var Tmpl *template.Template

// DB is the global database variable
var DB *sqlx.DB

// GinEngine is the global gin engine variable
var GinEngine *gin.Engine
