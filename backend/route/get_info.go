package route

import "github.com/gin-gonic/gin"

func GetInfo(c *gin.Context) {
	token := c.GetString("token")
	userId := c.GetString("user_id")
	c.JSON(200, gin.H{"token": token, "user_id": userId})
}
