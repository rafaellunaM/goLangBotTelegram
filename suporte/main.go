package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.ForwardedByClientIP = true
	r.SetTrustedProxies([]string{"127.0.0.1"})

	r.GET("/", func(c *gin.Context) {
		c.String(200, "Seja bem-vindo! Por favor, digite seu nome.")
	})

	r.GET("/pergunta", func(c *gin.Context) {
		c.String(200, "Por favor, digite o seu problema.")
	})

	r.GET("/aguarde", func(c *gin.Context) {
		c.String(200, "Por favor, aguarde 1 minuto para seu atendimento.")
	})

	if err := r.Run(":6060"); err != nil {
		log.Fatalf("Não foi possível iniciar o servidor: %v", err)
	}
}
