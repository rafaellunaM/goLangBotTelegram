package atendente

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"sync"

	"github.com/go-telegram/bot"

	"botTelegram/crud"
	dbConfig "botTelegram/dbconfig"
)

var AState = make(map[int64]string)
var mu sync.Mutex

var db *sql.DB
var err error
var orders []dbConfig.Order

func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}

var pedido = regexp.MustCompile(`^\d+$`)

var (
	root    = getEnv("ATENDENTE_API_ROOT", "http://localhost:5050")
	order   = getEnv("ATENDENTE_API_ORDER", "http://localhost:5050/order")
	aguarda = getEnv("ATENDENTE_API_AGUARDA", "http://localhost:5050/aguarda")
)

func HandleAttendant(ctx context.Context, b *bot.Bot, chatID int64) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   "Por favor, nos fale o seu pedido para que possamos atender com mais rapidez. Aguarde um momento que um de nossos atendentes irá lhe ajudar em breve.",
	})

	SetAState(chatID, "order-selected")
}

func SetAState(chatID int64, state string) {
	mu.Lock()
	defer mu.Unlock()
	AState[chatID] = state
}
func GetAState(chatID int64) string {
	mu.Lock()
	defer mu.Unlock()
	return AState[chatID]
}

func HandleOrder(ctx context.Context, b *bot.Bot, chatID int64, answer string) {
	request, err := http.Get(aguarda)
	if err != nil {
		log.Printf("Erro ao fazer requisição GET para %s: %v", aguarda, err)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "Erro ao acessar o endpoint de atendente.",
		})
		return
	}
	defer request.Body.Close()

	if pedido.MatchString(answer) {
		// Obtém os pedidos e produtos
		orders, err := crud.GetOrders(answer)
		// produtos, err := crud.GetProducts()
		if err != nil {

			fmt.Println(err)

			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: chatID,
				Text:   "Erro ao obter informações sobre pedidos.",
			})
			return
		}

		if orders != "" {

			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: chatID,
				Text:   "Seu pedido " + answer + " tem os seguintes produtos:\n" + orders,
			})
			SetAState(chatID, "")
		} else {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: chatID,
				Text:   "Seu pedido não consta na base de dados, digite novamente",
			})
		}

	} else {
		// Envia a resposta do usuário
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "Pedido inválido, digite novamente",
		})
	}
}

// func OrdersToString(orders []dbConfig.Order) string {
// 	var builder strings.Builder
// 	for _, order := range orders {
// 		// Construa a string para cada pedido
// 		builder.WriteString(fmt.Sprintf("ID do Pedido: %s, ID do Produto: %s\n", order.OrderID, order.ProductID))
// 	}
// 	return builder.String()
// }

// func FormatOrders(orders []dbConfig.Order, produtos []dbConfig.Article) string {
// 	var result string
// 	for _, order := range orders {
// 		var productPrice float32
// 		for _, produto := range produtos {
// 			if strconv.Itoa(produto.ID) == order.ProductID {
// 				productPrice = produto.Price
// 				break
// 			}
// 		}
// 		result += fmt.Sprintf("ID: %s, Produto: %s, Preço: %.2f\n", order.OrderID, order.ProductID, productPrice)
// 	}
// 	return result
// }

// func GetOrders(answer string) (string, error) {
// 	if err := acessDatabase(); err != nil {
// 		return "", fmt.Errorf("erro ao acessar o banco de dados: %w", err)

// 	}
// 	rows, err := db.Query("SELECT products.ID AS ProductID, products.Name AS ProductName, products.Price AS ProductPrice FROM orders INNER JOIN products ON orders.produto = products.ID;")
// 	if err != nil {
// 		return "", fmt.Errorf("erro ao executar a consulta: %w", err)
// 	}
// 	defer rows.Close()

// 	var result string
// 	for rows.Next() {
// 		var productID string
// 		var productName string
// 		var productPrice float64

// 		if err := rows.Scan(&productID, &productName, &productPrice); err != nil {
// 			return "", fmt.Errorf("erro ao ler os dados: %w", err)
// 		}
// 		result += fmt.Sprintf(" %s, | %s | , %.2f\n", productID, productName, productPrice)
// 	}

// 	if err := rows.Err(); err != nil {
// 		return "", fmt.Errorf("erro durante a iteração dos resultados: %w", err)
// 	}

// 	return result, nil
// }

func acessDatabase() error {
	fmt.Println("Acessando o banco de dados...")
	var err error
	db, err = sql.Open(dbConfig.PostgresDriver, dbConfig.DataSourceName)
	fmt.Println("Acessou...")
	if err != nil {
		return fmt.Errorf("falha ao abrir a conexão com o banco de dados: %w", err)
	}

	if err := db.Ping(); err != nil {
		return fmt.Errorf("falha ao conectar o banco de dados: %w", err)
	}
	return nil
}
