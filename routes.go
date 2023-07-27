package auth

import "github.com/gin-gonic/gin"

func RegisterRoutes(c *gin.Engine) {
	c.POST("/api/v1/login", Login)
}
