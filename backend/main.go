package main

import (
	"backend/route"

	"github.com/LyricTian/gzip"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.SecureJsonPrefix("")
	r.MaxMultipartMemory = (1 << 20) * 30
	r.Use(gzip.Gzip(gzip.DefaultCompression))
	r.Use(gin.Recovery())
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "X-Finger-Print"}
	r.Use(cors.New(corsConfig))
	route.Handler(r)
	r.Run(":8089")
}
