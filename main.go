package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"

	"botTelegram/crud"

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

func handlerResponse(ctx context.Context, b *bot.Bot, update *models.Update) {
	answer := update.Message.Text
	chatID := update.Message.Chat.ID

	switch answer {
	case "1":
		urlHello := "http://localhost:6060"
		urlPergunta := "http://localhost:6060/pergunta"
		urlAguarde := "http://localhost:6060/aguarde"

		hello, err := http.Get(urlHello)
		request, err := http.Get(urlPergunta)
		wait, err := http.Get(urlAguarde)

		if err != nil {
			log.Printf("Erro ao fazer requisição GET para %s: %v", urlHello, err)
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: chatID,
				Text:   "Erro ao acessar o endpoint de suporte.",
			})
			return
		}

		if err != nil {
			log.Printf("Erro ao fazer requisição GET para %s: %v", urlPergunta, err)
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: chatID,
				Text:   "Erro ao acessar o endpoint de suporte.",
			})
			return
		}

		if err != nil {
			log.Printf("Erro ao fazer requisição GET para %s: %v", urlAguarde, err)
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: chatID,
				Text:   "Erro ao acessar o endpoint de suporte.",
			})
			return
		}

		defer hello.Body.Close()
		defer request.Body.Close()
		defer wait.Body.Close()

		helloBody, err := ioutil.ReadAll(hello.Body)
		if err != nil {
			log.Printf("Erro ao ler a resposta do endpoint: %v", err)
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: chatID,
				Text:   "Erro ao ler a resposta do endpoint de suporte.",
			})
			return
		}

		messageHello := fmt.Sprintf("sd_bot: %s", string(helloBody))

		_, err = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    chatID,
			Text:      messageHello,
			ParseMode: "HTML",
		})

		if err != nil {
			log.Printf("Erro ao enviar mensagem: %v", err)
		}

		requestBody, err := ioutil.ReadAll(request.Body)
		messageRequest := fmt.Sprintf("sd_bot: %s", string(requestBody))

		waitBody, err := ioutil.ReadAll(wait.Body)
		messageWait := fmt.Sprintf("sd_bot: %s", string(waitBody))

	case "2":
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

	case "3":
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

func processUpdate(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message != nil {
		chatID := update.Message.Chat.ID
		if update.Message.Text == "1" || update.Message.Text == "2" || update.Message.Text == "3" {
			handlerResponse(ctx, b, update)
		} else {
			sendWelcomeMessage(ctx, b, chatID)
		}
	}
}
