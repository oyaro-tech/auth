package auth

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(c *gin.Engine) {
	c.POST("/login", Login)
	c.POST("/logout", TokenAuthMiddleware, Logout)
}
