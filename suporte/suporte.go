package suporte

import (
	"botTelegram/crud"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"

	"github.com/go-telegram/bot"
)

var userStates = make(map[int64]string)
var mu sync.Mutex

func HanlderHelloUser(ctx context.Context, b *bot.Bot, chatID int64) {
	urlHello := "http://localhost:6060"
	hello, err := http.Get(urlHello)
	if err != nil {
		log.Printf("Erro ao fazer requisição GET para %s: %v", urlHello, err)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "Erro ao acessar o endpoint de suporte.",
		})
		return
	}
	defer hello.Body.Close()

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
	SetUserState(chatID, "awaiting_username")
}

func HandlerUserName(ctx context.Context, b *bot.Bot, chatID int64, username string) {
	urlPergunta := "http://localhost:6060/pergunta"
	request, err := http.Get(urlPergunta)
	if err != nil {
		log.Printf("Erro ao fazer requisição GET para %s: %v", urlPergunta, err)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "Erro ao acessar o endpoint de suporte.",
		})
		return
	}
	defer request.Body.Close()

	requestBody, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.Printf("Erro ao ler a resposta do endpoint: %v", err)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "Erro ao ler a resposta do endpoint de suporte.",
		})
		return
	}

	messageRequest := fmt.Sprintf("sd_bot: %s", string(requestBody))

	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    chatID,
		Text:      messageRequest,
		ParseMode: "HTML",
	})
	if err != nil {
		log.Printf("Erro ao enviar mensagem: %v", err)
	}
	SetUserState(chatID, "awaiting_issues")
}

func HandlerIssues(ctx context.Context, b *bot.Bot, chatID int64) {
	urlAguarde := "http://localhost:6060/aguarde"
	wait, err := http.Get(urlAguarde)
	if err != nil {
		log.Printf("Erro ao fazer requisição GET para %s: %v", urlAguarde, err)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "Erro ao acessar o endpoint de suporte.",
		})
		return
	}
	defer wait.Body.Close()

	waitBody, err := ioutil.ReadAll(wait.Body)
	if err != nil {
		log.Printf("Erro ao ler a resposta do endpoint: %v", err)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "Erro ao ler a resposta do endpoint de suporte.",
		})
		return
	}

	messageWait := fmt.Sprintf("sd_bot: %s", string(waitBody))

	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    chatID,
		Text:      messageWait,
		ParseMode: "HTML",
	})
	if err != nil {
		log.Printf("Erro ao enviar mensagem: %v", err)
	}
	SetUserState(chatID, "complete")
}

func HandleProducts(ctx context.Context, b *bot.Bot, chatID int64) {
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
}

func HandleAttendant(ctx context.Context, b *bot.Bot, chatID int64) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   "Por favor, nos fale o seu pedido para que possamos atender com mais rapidez. Aguarde um momento que um de nossos atendentes irá lhe ajudar em breve.",
	})
}

func SendWelcomeMessage(ctx context.Context, b *bot.Bot, chatID int64) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   welcome(),
	})
}

func GetUserStates(chatID int64) string {
	mu.Lock()
	defer mu.Unlock()
	return userStates[chatID]
}

func SetUserState(chatID int64, state string) {
	mu.Lock()
	defer mu.Unlock()
	userStates[chatID] = state
}

func welcome() string {
	return "Olá, seja bem-vindo a lombratec\nEscolha uma das opções:\n1 - Suporte\n2 - Produtos\n3 - Atendente"
}
