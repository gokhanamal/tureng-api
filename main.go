package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gokhanamal/tureng-api/controller"
)

func main() {
	router := gin.Default()
	router.GET("/translate", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"error":    "Missing query, you should add to your request phrase that want to translate",
			"response": "",
		})
	})

	router.GET("/translate/:phrase", func(c *gin.Context) {
		phrase := c.Param("phrase")
		if phrase == "" {
			c.JSON(200, gin.H{
				"error":    "Missing query, you should add to your request phrase that want to translate",
				"response": "",
			})
			return
		}
		response, err := controller.FetchFromTureng(phrase)
		fmt.Println(response)
		if err != nil {
			c.JSON(200, gin.H{
				"error":    "Someting went wrong while fetching the phrases from Tureng",
				"response": "",
			})
			return
		}
		c.JSON(200, gin.H{
			"error":    "",
			"response": response,
		})
	})

	err := router.Run(":8080")

	if err != nil {
		fmt.Printf("%s", err)
	}
}
