package server

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
)

func getGeoJsonLineSingle(c *gin.Context) {

	dsId := c.Query("data-source-filepath")
	if dsId == "" {
		c.JSON(400, gin.H{
			"message": "data-source-filepath query param is required",
		})
		return
	}

	ctx := context.Background()
	val, err := q.GetLineByFilepath(ctx, dsId)
	if err != nil {
		log.Printf("error while fetching data sources: %s", err)
		c.JSON(503, gin.H{
			"message": "error",
		})
		return
	}

	c.Data(200, "application/geo+json", []byte(val))
}

func init() {
	engine.GET("/data-sources/lines", getGeoJsonLineSingle)
}
