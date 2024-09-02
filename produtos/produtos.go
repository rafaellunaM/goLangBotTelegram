package produtos

import (
	"botTelegram/crud"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"sync"

	"github.com/go-telegram/bot"
)

var userStates = make(map[int64]string)
var mu sync.Mutex

var nameRegex = regexp.MustCompile(`^[A-Za-z\s]+$`)
var numberRegex = regexp.MustCompile(`^\d+$`)
var cpfRegex = regexp.MustCompile(`^\d+$`)
var re = regexp.MustCompile(`(?i)^sim$`)

func AnswerTreatment(answer string) bool {

	if (re.MatchString(answer)) {
		return true
	}
	return false
}

func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}

var (
	root     = getEnv("SUPORTE_API_ROOT", "http://localhost:7070")
	produtos = getEnv("SUPORTE_API_PRODUTOS", "http://localhost:7070/produtos")
	produto  = getEnv("SUPORTE_API_PRODUTO", "http://localhost:7070/produto")
)


func NameTratment(awnser string) bool {

	if nameRegex.MatchString(awnser) {
		log.Printf("Mensagem recebida: %s", awnser)
		return true
	} else {
		log.Printf("Mensagem inválida recebida: %s", awnser)
	}
	return false

}

func ProdutosTratment(produtos string) bool {

	if numberRegex.MatchString(produtos) {
		log.Printf("Mensagem recebida: %s", produtos)
		return true
	} else {
		log.Printf("Mensagem inválida recebida: %s", produtos)
	}
	return false

}

func ProdutoTratment(produto string) bool {

	if cpfRegex.MatchString(produto) {
		log.Printf("Mensagem recebida: %s", produto)
		return true
	} else {
		log.Printf("Mensagem inválida recebida: %s", produto)
	}
	return false

}

func HanlderHelloUser(ctx context.Context, b *bot.Bot, chatID int64) {
	hello, err := http.Get(root)
	if err != nil {
		log.Printf("Erro ao fazer requisição GET para %s: %v", root, err)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "Erro ao acessar o endpoint de produtos.",
		})
		return
	}
	defer hello.Body.Close()

	helloBody, err := ioutil.ReadAll(hello.Body)
	if err != nil {
		log.Printf("Erro ao ler a resposta do endpoint: %v", err)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "Erro ao ler a resposta do endpoint de produto.",
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
	SetUserState(chatID, "awaiting_answer")
}

func HandlerProdutos(ctx context.Context, b *bot.Bot, chatID int64, produtos string) {

	request, err := http.Get(produtos)
	if err != nil {
		log.Printf("Erro ao fazer requisição GET para %s: %v", produtos, err)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "Erro ao acessar o endpoint de produto.",
		})
		return
	}
	defer request.Body.Close()

	requestBody, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.Printf("Erro ao ler a resposta do endpoint: %v", err)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "Erro ao ler a resposta do endpoint de produtos.",
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

	SetUserState(chatID, "awaiting_produto")
}
func HandlerProduto(ctx context.Context, b *bot.Bot, chatID int64, produto string) {
	request, err := http.Get(produto)
	if err != nil {
		log.Printf("Erro ao fazer requisição GET para %s: %v", produto, err)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "Erro ao acessar o endpoint de produtos.",
		})
		return
	}
	defer request.Body.Close()

	requestBody, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.Printf("Erro ao ler a resposta do endpoint: %v", err)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "Erro ao ler a resposta do endpoint de produtos.",
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
	SetUserState(chatID, "")
}
/* func HandlerUserPhone(ctx context.Context, b *bot.Bot, chatID int64, phone string) {
	request, err := http.Get(pergunta)
	if err != nil {
		log.Printf("Erro ao fazer requisição GET para %s: %v", pergunta, err)
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
} */

/* func HandlerIssues(ctx context.Context, b *bot.Bot, chatID int64) {
	wait, err := http.Get(aguarde)
	if err != nil {
		log.Printf("Erro ao fazer requisição GET para %s: %v", aguarde, err)
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
} */

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
