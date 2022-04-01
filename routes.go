package auth

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(c *gin.Engine) {
	c.POST("/auth/login", Login)
	c.POST("/auth/logout", Logout)
}
