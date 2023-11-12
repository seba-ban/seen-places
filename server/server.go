package server

import (
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
	TemplatesExt  string
}

func (c *ServerConfig) Run() {
	ctx, dbConn, err := commonflags.DbConnConfig.OpenDbConnection()
	utils.PanicOnError(err)
	defer dbConn.Close(*ctx)

	q = queries.New(dbConn)

	engine.LoadHTMLGlob(path.Join(c.TemplatesPath, "*"+c.TemplatesExt))

	// TODO: move to config
	engine.Run() // listen and serve on 0.0.0.0:8080
}
