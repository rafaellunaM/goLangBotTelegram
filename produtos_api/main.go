package main

import (
	"github.com/gin-gonic/gin"

	"log"
)

func main() {

	r := gin.Default()

	r.ForwardedByClientIP = true
	r.SetTrustedProxies([]string{"127.0.0.1"})

	r.GET("/", func(c *gin.Context) {
		c.String(200, "Você gostaria de ver os nossos produtos? digite sim ou não")
	})

	r.GET("/produtos", func(c *gin.Context) {
		c.String(200, "Você deseja consultar todos os produtos disponíveis?")
	})

	r.GET("/produto", func(c *gin.Context) {
		c.String(200, "Você deseja consultar todos os produtos disponíveis?")
	})

	if err := r.Run(":7070"); err != nil {
		log.Fatalf("Não foi possível iniciar o servidor: %v", err)
	}
}
