package main

import (
	"encoding/json"
	"fmt"
	"github.com/dan-lugg/go-commands-example/app"
	"github.com/dan-lugg/go-commands-example/app/util"
	"github.com/dan-lugg/go-commands/commands"
	"github.com/dan-lugg/go-commands/openapi"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

func main() {
	err, container := app.BuildContainer()
	if err != nil {
		fmt.Printf("Failed to build container: %v\n", err)
		return
	}

	router := gin.Default()
	router.GET("/docs", func(c *gin.Context) {
		mappingCatalog := container.Get(util.TypeNameFor[commands.MappingCatalog]()).(*commands.MappingCatalog)
		handlerCatalog := container.Get(util.TypeNameFor[commands.HandlerCatalog]()).(*commands.HandlerCatalog)

		specWriter := openapi.NewSpecWriter(mappingCatalog, handlerCatalog)
		c.Header("Content-Type", "application/json")
		err := specWriter.WriteSpec(c.Writer)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("error writing OpenAPI spec: %v", err)})
			return
		}

		c.Status(http.StatusOK)
	})
	router.POST("/commands/:name", func(c *gin.Context) {
		manager := container.Get(util.TypeNameFor[commands.Manager]()).(*commands.Manager)
		reqName := c.Param("name")
		reqData, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("error reading request body: %v", err)})
			return
		}

		res, err := manager.Handle(reqName, reqData, c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("error handling command: %v", err)})
			return
		}

		resData, err := json.Marshal(res)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("error encoding response: %v", err)})
			return
		}

		c.Data(http.StatusOK, "application/json", resData)
	})

	err = router.Run(":8080")
	if err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
		return
	}
}
