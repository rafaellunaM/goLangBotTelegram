package main
import (
	"github.com/gin-gonic/gin"
	"botTelegram/crud"
	"net/http"
	//"context"
	"fmt"
	"log"
	//"os"
	//"os/signal"

	//"botTelegram/crud"
	//"botTelegram/produtos"

	//"github.com/go-telegram/bot"
	//"github.com/go-telegram/bot/models"
	//"github.com/joho/godotenv"
)

func main() {

	r := gin.Default()

	r.ForwardedByClientIP = true
	r.SetTrustedProxies([]string{"127.0.0.1"})

	r.GET("/", func(c *gin.Context) {
		c.String(200, "Nesta API você pode consultar nossos produtos disponíveis.")
	})

	r.GET("/produtos", func(c *gin.Context) {
		products, err := crud.GetProducts()
		if err != nil {
			log.Printf("Erro ao buscar produtos: %v", err)
			c.String(http.StatusInternalServerError, "Erro ao buscar produtos.")
			return
		}

		var productList string
		for _, product := range products {
			productList += fmt.Sprintf("%s: R$%.2f\n", product.Name, product.Price)
		}

		c.String(http.StatusOK, "Lista de Produtos:\n"+productList)
	})

	if err := r.Run(":7070"); err != nil {
		log.Fatalf("Não foi possível iniciar o servidor: %v", err)
	}
}