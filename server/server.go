package server

import (
	"html/template"
	"path"

	"github.com/gin-gonic/gin"
	commonflags "github.com/seba-ban/seen-places/commonFlags"
	"github.com/seba-ban/seen-places/queries"
	"github.com/seba-ban/seen-places/utils"
)

var engine = gin.Default()
var q *queries.Queries

type ServerConfig struct {
	TemplatesPath string
}

func (c *ServerConfig) Run() {
	dbConn, err := commonflags.DbConnConfig.OpenDbConnection()
	utils.PanicOnError(err)
	defer dbConn.Close()

	q = queries.New(dbConn)

	engine.SetFuncMap(template.FuncMap{
		"formatTime":         utils.FormatTime,
		"formatTimeDuration": utils.FormatTimeDuration,
	})
	engine.LoadHTMLGlob(path.Join(c.TemplatesPath, "*"))
	// TODO: move to config
	engine.Run() // listen and serve on 0.0.0.0:8080
}
