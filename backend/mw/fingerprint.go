package mw

import (
	"backend/enigma_im"

	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		fp := c.GetHeader("X-Finger-Print")
		if fp == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "missing fingerprint"})
			return
		}
		token := enigma_im.GetUserToken(c.Request.Context(), fp)
		if token == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "invalid fingerprint"})
			return
		}
		c.Set("user_id", fp)
		c.Set("token", token)
		c.Next()
	}
}
