package route

import (
	"backend/enigma_im"
	"log"

	"github.com/gin-gonic/gin"
)

func SendChannelMsg(c *gin.Context) {
	var (
		ctx    = c.Request.Context()
		userId = c.GetString("user_id")
		req    = struct {
			RoomId  string `json:"room_id" binding:"required"`
			Content string `json:"content" binding:"required"`
		}{}
	)
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("invalid request: %v", err)
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	if err := enigma_im.SendChannelMsg(ctx, userId, req.RoomId, req.Content); err != nil {
		log.Printf("failed to send channel msg: %v", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": "ok"})
}
