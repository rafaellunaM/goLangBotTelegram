package main

import (
	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()
	r.ForwardedByClientIP = true
	r.SetTrustedProxies([]string{"127.0.0.1"})

	r.GET("/", func(c *gin.Context) {
		c.String(200, "oi")
	})

	r.Run(":6060")
}
