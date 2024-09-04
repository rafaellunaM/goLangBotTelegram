package atendente

import (
	"botTelegram/suporte"
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/go-telegram/bot"

	dbConfig "botTelegram/dbconfig"
)

var userStates = make(map[int64]string)
var mu sync.Mutex

var db *sql.DB
var err error
var orders []dbConfig.Order

func HandleAttendant(ctx context.Context, b *bot.Bot, chatID int64) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   "Por favor, nos fale o seu pedido para que possamos atender com mais rapidez. Aguarde um momento que um de nossos atendentes irá lhe ajudar em breve.",
	})

	suporte.SetUserState(chatID, "order-selected")
}

// func HandleOrder(ctx context.Context, b *bot.Bot, chatID int64, answer string) {
// 	urlPergunta := "http://localhost:5050/aguarda"
// 	request, err := http.Get(urlPergunta)
// 	b.SendMessage(ctx, &bot.SendMessageParams{
// 		ChatID: chatID,
// 		Text:   "Seu pedido é " + answer,
// 	})

// 	GetOrders()

// 	if err != nil {
// 		log.Printf("Erro ao fazer requisição GET para %s: %v", urlPergunta, err)
// 		b.SendMessage(ctx, &bot.SendMessageParams{
// 			ChatID: chatID,
// 			Text:   "Erro ao acessar o endpoint de atendente.",
// 		})
// 		return
// 	}

// 	suporte.SetUserState(chatID, "")

// 	defer request.Body.Close()
// }

func HandleOrder(ctx context.Context, b *bot.Bot, chatID int64, answer string) {
	urlPergunta := "http://localhost:5050/aguarda"
	request, err := http.Get(urlPergunta)
	if err != nil {
		log.Printf("Erro ao fazer requisição GET para %s: %v", urlPergunta, err)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "Erro ao acessar o endpoint de atendente.",
		})
		suporte.SetUserState(chatID, " ")
		return
	}
	defer request.Body.Close()

	// Envia a resposta do usuário
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   "Seu pedido é " + answer,
	})

	// Chama GetOrders para obter a lista de pedidos
	err = GetOrders()
	if err != nil {
		log.Fatalf("Erro ao obter pedidos: %v", err)
	}

	if err != nil {
		log.Printf("Erro ao obter pedidos: %v", err)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "Erro ao acessar os pedidos.",
		})
		return
	}

	// Formata e envia os pedidos para o usuário
	if len(orders) > 0 {
		for _, order := range orders {
			message := fmt.Sprintf("Pedido ID: %d, Produto ID: %d", order.OrderID, order.ProductID)
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: chatID,
				Text:   message,
			})
		}
	} else {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "Nenhum pedido encontrado.",
		})
	}

}

func GetOrders() error {
	if err := acessDatabase(); err != nil {
		return fmt.Errorf("erro ao acessar o banco de dados: %w", err)
	}

	rows, err := db.Query("SELECT * FROM OrderProducts")
	if err != nil {
		return fmt.Errorf("erro ao executar a consulta: %w", err)
	}
	defer rows.Close()

	fmt.Println("Resultados da consulta:")
	for rows.Next() {
		var orderID int
		var productID int
		if err := rows.Scan(&orderID, &productID); err != nil {
			return fmt.Errorf("erro ao ler os dados: %w", err)
		}
		fmt.Printf("OrderID: %d, ProductID: %d\n", orderID, productID)
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("erro durante a iteração dos resultados: %w", err)
	}
	return nil
}

// func GetOrders() ([]dbConfig.Order, error) {
// 	if err := acessDatabase(); err != nil {
// 		return nil, fmt.Errorf("erro ao acessar o banco de dados: %w", err)
// 	}

// 	rows, err := db.Query("SELECT OrderID, ProductID FROM OrderProducts")
// 	if err != nil {
// 		return nil, fmt.Errorf("erro ao executar a consulta: %w", err)
// 	}
// 	defer rows.Close()

// 	var orders []dbConfig.Order
// 	for rows.Next() {
// 		var order dbConfig.Order
// 		if err := rows.Scan(&order.OrderID, &order.ProductID); err != nil {
// 			return nil, fmt.Errorf("erro ao ler os dados: %w", err)
// 		}
// 		orders = append(orders, order)
// 	}

// 	if err := rows.Err(); err != nil {
// 		return nil, fmt.Errorf("erro durante a iteração dos resultados: %w", err)
// 	}
// 	return orders, nil
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

func SetUserState(chatID int64, state string) {
	mu.Lock()
	defer mu.Unlock()
	userStates[chatID] = state
}
