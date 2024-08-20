package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	"botTelegram/crud"
	"botTelegram/suporte"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/joho/godotenv"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		bot.WithDefaultHandler(processUpdate),
	}

	if err := godotenv.Load(); err != nil {
		log.Fatalf("Erro ao carregar o arquivo .env: %v", err)
	}

	botToken := os.Getenv("toke_telegram")
	if botToken == "" {
		log.Fatal("A variável de ambiente 'toke_telegram' não está definido")
	}

	b, err := bot.New(botToken, opts...)
	if err != nil {
		panic(err)
	}

	b.Start(ctx)
}

func welcome() string {
	return "Olá, seja bem-vindo a lombratec\nEscolha uma das opções:\n1 - Suporte\n2 - Produtos\n3 - Atendente"
}

func sendWelcomeMessage(ctx context.Context, b *bot.Bot, chatID int64) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   welcome(),
	})
}

func processUpdate(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message != nil {
		handlerResponse(ctx, b, update)
	}
}

func handlerResponse(ctx context.Context, b *bot.Bot, update *models.Update) {
	answer := update.Message.Text
	chatID := update.Message.Chat.ID

	state := suporte.GetUserStates(chatID)

	switch {
	case answer == "1" && state == "":
		suporte.HanlderHelloUser(ctx, b, chatID)
	case state == "awaiting_username":
		suporte.HandlerUserName(ctx, b, chatID, answer)
	case state == "awaiting_issues":
		suporte.HandlerIssues(ctx, b, chatID)
	case answer == "2":
		products, err := crud.GetProducts()
		if err != nil {
			log.Printf("Erro ao buscar produtos: %v", err)
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: chatID,
				Text:   "Erro ao buscar produtos.",
			})
			return
		}

		var productList string
		for _, product := range products {
			productList += fmt.Sprintf("%s: R$%.2f\n", product.Name, product.Price)
		}

		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "Lista de Produtos:\n" + productList,
		})

	case answer == "3":
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "Por favor, nos fale o seu pedido para que possamos atender com mais rapidez. Aguarde um momento que um de nossos atendentes irá lhe ajudar em breve.",
		})

	default:
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "Escolha uma das opções:\n1 - Suporte\n2 - Produtos\n3 - Atendente",
		})
	}
}
