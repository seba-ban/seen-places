package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func rootHandler(c *gin.Context) {
	tmpl := TemplateIndex{}
	c.HTML(http.StatusOK, tmpl.TemplateName(), tmpl.TemplateInputMap())
}

func init() {
	engine.GET("/", rootHandler)
}
