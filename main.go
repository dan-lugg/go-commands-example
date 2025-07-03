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
		mappingRegistry := container.Get(util.TypeNameFor[commands.MappingRegistry]()).(*commands.MappingRegistry)
		handlerRegistry := container.Get(util.TypeNameFor[commands.HandlerRegistry]()).(*commands.HandlerRegistry)

		specWriter := openapi.NewSpecWriter(mappingRegistry, handlerRegistry)
		c.Header("Content-Type", "application/json")
		err = specWriter.WriteSpec(c.Writer)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("error writing OpenAPI spec: %v", err)})
			return
		}

		c.Status(http.StatusOK)
	})
	router.POST("/commands/:name", func(c *gin.Context) {
		mappingRegistry := container.Get(util.TypeNameFor[commands.MappingRegistry]()).(*commands.MappingRegistry)
		decoderRegistry := container.Get(util.TypeNameFor[commands.DecoderRegistry]()).(*commands.DecoderRegistry)
		handlerRegistry := container.Get(util.TypeNameFor[commands.HandlerRegistry]()).(*commands.HandlerRegistry)

		reqName := c.Param("name")
		reqData, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("error reading request body: %v", err)})
			return
		}

		reqType, err := mappingRegistry.ByName(reqName)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("error finding mapped type for request: %v", err)})
			return
		}

		req, err := decoderRegistry.Decode(reqType, reqData)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("error decoding request: %v", err)})
			return
		}

		res, err := handlerRegistry.Handle(req, c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("error handling request: %v", err)})
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
