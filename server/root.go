package server

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func rootHandler(c *gin.Context) {

	ctx := context.Background()
	geo, err := q.GetGeoJson(ctx)
	if err != nil {
		log.Printf("error while fetching geojson: %s", err)
		// TODO: return html
		c.JSON(503, gin.H{
			"message": "error",
		})
		return
	}

	tmpl := TemplateIndex{
		GeoJSON: &geo,
	}
	c.HTML(http.StatusOK, tmpl.TemplateName(), tmpl.TemplateInputMap())
}

func init() {
	engine.GET("/", rootHandler)
}
