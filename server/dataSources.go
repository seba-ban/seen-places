package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"log"

	"github.com/gin-gonic/gin"
	"github.com/seba-ban/seen-places/queries"
	"github.com/seba-ban/seen-places/utils"
)

type DataSource struct {
	Type      string    `json:"type"`
	Filename  string    `json:"filename"`
	Size      int64     `json:"size"`
	CreatedAt time.Time `json:"createdAt"`
}

func getDataSourcesHandler(c *gin.Context) {
	qc := utils.QueryConverter{Ctx: c}

	params := &queries.GetDataSourcesParams{}

	originalFilename, err := qc.GetPgtypeText("filename")
	if err == nil {
		params.OriginalFilename = *originalFilename
	}

	typeQ, err := qc.GetPgtypeText("type")
	if err == nil {
		params.Type = *typeQ
	}

	startBefore, err := qc.GetPgtypeTimestamp("startBefore")
	if errors.Is(err, utils.ErrorQueryNotValid) {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "startBefore is not a valid RFC3339 time",
		})
		return
	} else if err == nil {
		params.StartBefore = *startBefore
	}

	startAfter, err := qc.GetPgtypeTimestamp("startAfter")
	if errors.Is(err, utils.ErrorQueryNotValid) {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "startAfter is not a valid RFC3339 time",
		})
		return
	} else if err == nil {
		params.StartAfter = *startAfter
	}

	ctx := context.Background()

	ds, err := q.GetDataSources(ctx, params)
	if err != nil {
		log.Printf("error while fetching data sources: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "error while fetching data sources",
		})
		return
	}

	// to avoid returning null in the response
	if len(ds) == 0 {
		ds = []queries.GetDataSourcesRow{}
	}

	c.JSON(http.StatusOK, gin.H{
		"data": ds,
	})

}

func getDataSourceByPolygon(c *gin.Context) {
	qc := utils.QueryConverter{Ctx: c}

	points := make([]float64, 8)

	for i, queryName := range []string{"x1", "y1", "x2", "y2", "x3", "y3", "x4", "y4"} {
		point, err := qc.GetFloat(queryName)

		if err != nil {
			switch err {
			case utils.ErrorQueryNotValid:
				c.JSON(http.StatusBadRequest, gin.H{
					"message": fmt.Sprintf("query %s is not a valid float", queryName),
				})
				return
			case utils.ErrorQueryNotDefined:
				c.JSON(http.StatusBadRequest, gin.H{
					"message": fmt.Sprintf("query %s is required", queryName),
				})
				return
			default:
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": fmt.Sprintf("error while parsing %s", queryName),
				})
				return
			}
		}

		points[i] = point
	}

	linestring := fmt.Sprintf("LINESTRING(%f %f, %f %f, %f %f, %f %f, %f %f)", points[0], points[1], points[2], points[3], points[4], points[5], points[6], points[7], points[0], points[1])

	ctx := context.Background()

	ds, err := q.GetDataSourcesFromPolygon(ctx, linestring)
	if err != nil {
		log.Printf("error while fetching data sources: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "error while fetching data sources",
		})
		return
	}

	// to avoid returning null in the response
	if len(ds) == 0 {
		ds = []queries.GetDataSourcesFromPolygonRow{}
	}

	if utils.IsHTMXRequest(c) {
		t := TemplateDataSourceEls{DataSources: ds}
		c.HTML(http.StatusOK, t.TemplateName(), t.TemplateInputMap())
	} else {
		c.JSON(http.StatusOK, gin.H{
			"data": ds,
		})
	}

}

func init() {
	engine.GET("/data-sources", getDataSourcesHandler)
	engine.GET("/data-sources/within-polygon", getDataSourceByPolygon)
}
