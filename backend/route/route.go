package route

import (
	"backend/mw"

	"github.com/gin-gonic/gin"
)

func Handler(e *gin.Engine) {
	v1 := e.Group("/chat_demo/v1")
	v1.Use(mw.Auth())
	v1.GET("get_info", GetInfo)
	v1.POST("send_channel_msg", SendChannelMsg)
}
