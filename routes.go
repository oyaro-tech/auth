package auth

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(c *gin.Engine) {
	c.POST("/register", Register)
	c.POST("/login", Login)
	c.POST("/logout", TokenAuthMiddleware, Logout)
}
