package api

import "github.com/gin-gonic/gin"

// HelloHandler handles hello
func HelloHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"text": "Hello World.",
	})
}
