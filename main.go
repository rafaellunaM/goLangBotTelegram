package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"

	"botTelegram/crud"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/joho/godotenv"
)

var userStates = make(map[int64]string)
var mu sync.Mutex

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
		update.Message.Text = ""

		hello, err := http.Get(urlHello)
		request, err := http.Get(urlPergunta)
		wait, err := http.Get(urlAguarde)

		//state := getUserStates(chatID)

		if err != nil {
			log.Printf("Erro ao fazer requisição GET para %s: %v", urlHello, err)
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: chatID,
				Text:   "Erro ao acessar o endpoint de suporte.",
			})
			return
		}

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

		setUserState(chatID, "awaiting_username")
		defer hello.Body.Close()

		if err != nil {
			log.Printf("Erro ao fazer requisição GET para %s: %v", urlPergunta, err)
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: chatID,
				Text:   "Erro ao acessar o endpoint de suporte.",
			})
			return
		}
		if userStates[chatID] == "awaiting_username" {
			requestBody, err := ioutil.ReadAll(request.Body)
			messageRequest := fmt.Sprintf("sd_bot: %s", string(requestBody))

			_, err = b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID:    chatID,
				Text:      messageRequest,
				ParseMode: "HTML",
			})

			if err != nil {
				log.Printf("Erro ao enviar mensagem: %v", err)
			}
			defer request.Body.Close()

			if err != nil {
				log.Printf("Erro ao fazer requisição GET para %s: %v", urlAguarde, err)
				b.SendMessage(ctx, &bot.SendMessageParams{
					ChatID: chatID,
					Text:   "Erro ao acessar o endpoint de suporte.",
				})
				return
			}

		}
		setUserState(chatID, "awaiting_issues")
		if userStates[chatID] == "awaiting_issues" {
			waitBody, err := ioutil.ReadAll(wait.Body)
			messageWait := fmt.Sprintf("sd_bot: %s", string(waitBody))

			_, err = b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID:    chatID,
				Text:      messageWait,
				ParseMode: "HTML",
			})
			if err != nil {
				log.Printf("Erro ao enviar mensagem: %v", err)
			}
			defer wait.Body.Close()

		}

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

func setUserState(chatID int64, state string) {
	mu.Lock()
	defer mu.Unlock()
	userStates[chatID] = state
}

func getUserStates(chatID int64) string {
	mu.Lock()
	defer mu.Unlock()
	return userStates[chatID]
}
