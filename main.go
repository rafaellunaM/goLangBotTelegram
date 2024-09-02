package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	"botTelegram/crud"
	"botTelegram/produtos"
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

type UserData struct {
	Name   string
	CPF    string
	Phone  string
	Issues string
}

var userResponses = make(map[int64]UserData)

func handlerResponse(ctx context.Context, b *bot.Bot, update *models.Update) {
	answer := update.Message.Text
	chatID := update.Message.Chat.ID

	if _, ok := userResponses[chatID]; !ok {
		userResponses[chatID] = UserData{}
	}

	state := suporte.GetUserStates(chatID)
	teste := produtos.GetTest(chatID)

	switch {
	case answer == "1" && state == "":
		suporte.HanlderHelloUser(ctx, b, chatID)
	case state == "awaiting_username" && suporte.NameTratment(answer) != false:
		userResponses[chatID] = UserData{Name: answer, CPF: userResponses[chatID].CPF, Phone: userResponses[chatID].Phone}
		suporte.HandlerUserName(ctx, b, chatID, answer)
	case state == "awaiting_cpf" && suporte.CpfTratmnet(answer) != false:
		userResponses[chatID] = UserData{Name: userResponses[chatID].Name, CPF: answer, Phone: userResponses[chatID].Phone}
		suporte.HandlerUserCpf(ctx, b, chatID, answer)
	case state == "awaiting_phone" && suporte.PhoneTratmnet(answer) != false:
		userResponses[chatID] = UserData{Name: userResponses[chatID].Name, CPF: userResponses[chatID].CPF, Phone: answer}
		suporte.HandlerUserPhone(ctx, b, chatID, answer)
	case state == "awaiting_issues":
		userResponses[chatID] = UserData{Name: userResponses[chatID].Name, CPF: userResponses[chatID].CPF, Phone: userResponses[chatID].Phone, Issues: answer}
		suporte.HandlerIssues(ctx, b, chatID)
		userData := userResponses[chatID]
		log.Printf("Dados do usuário %d: Nome= %s, CPF= %s, Telefone= %s, Problema= %s", chatID, userData.Name, userData.CPF, userData.Phone, userData.Issues)

		crud.SetUsers(userResponses[chatID].CPF, userResponses[chatID].Name, userResponses[chatID].Phone, userResponses[chatID].Issues)

	case answer == "2" && teste == "":
		produtos.HanlderHelloUser(ctx, b, chatID)
	case teste == "awaiting_answer":

		products, err := crud.GetProducts()
		if err != nil {
			log.Printf("Erro ao buscar produtos: %v", err)
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "Erro ao buscar produtos.",
			})
			return
		}

		var productList string
		for _, product := range products {
			productList += fmt.Sprintf("%s: R$%.2f\n", product.Name, product.Price)
		}

		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Lista de Produtos:\n" + productList,
		})

	case state == "awaiting_aluno" && produtos.NameTratment(answer) == false:
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "int te fode",
		})
		produtos.SetUserTest(chatID, "")

	case answer == "3":
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "Por favor, nos fale o seu pedido para que possamos atender com mais rapidez. Aguarde um momento que um de nossos atendentes irá lhe ajudar em breve.",
		})

	default:
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "Escolha uma das opções:\n1 - Suporte\n2 - Produtos\n3 - Atendente\n ou digite a resposta da pergunta anterior corretamente",
		})
	}
}
