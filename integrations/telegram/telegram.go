package telegram

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/da4nik/shopwatcher/db"
	"github.com/da4nik/shopwatcher/types"

	"gopkg.in/telegram-bot-api.v4"
)

var (
	outputChan chan types.Product
	bot        *tgbotapi.BotAPI
	inputChan  = make(chan types.Product, 10)
	done       = make(chan bool, 1)
	connected  = false
	commands   = map[string]func(*tgbotapi.Message){
		"/add":   addURL,
		"/list":  list,
		"/rm":    rm,
		"/start": help,
	}
)

// Start starts telegram bot
func Start(ctx context.Context, outChan chan types.Product) chan types.Product {
	outputChan = outChan

	go processBot(ctx)
	go listen()

	return inputChan
}

func processBot(ctx context.Context) {
	log := logger()

	var err error
	bot, err = tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_ACCESS_TOKEN"))
	if err != nil {
		log.Errorf("Unable to create new bot. %s", err.Error())
		return
	}

	// bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, _ := bot.GetUpdatesChan(u)

	connected = true

	for update := range updates {
		if update.Message == nil {
			continue
		}

		for command, funct := range commands {
			text := update.Message.Text
			if strings.HasPrefix(strings.ToLower(text), command) {
				funct(update.Message)
			}
		}
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
	}
}

func listen() {
	for {
		select {
		case product := <-inputChan:
			sendChanges(product)
		case <-done:
			return
		}
	}
}

func addURL(msg *tgbotapi.Message) {
	text := strings.ToLower(msg.Text)
	url := strings.Split(text, "/add ")[1]

	msgText := fmt.Sprintf("Url: %s added to watch list.", url)

	message := tgbotapi.NewMessage(msg.Chat.ID, msgText)
	message.DisableWebPagePreview = true
	bot.Send(message)

	product := types.Product{
		URL: strings.Trim(url, " \n"),
		Users: []types.User{types.User{
			ChatType: "telegram",
			ChatID:   strconv.FormatInt(msg.Chat.ID, 10),
		}},
	}
	outputChan <- product
}

func sendChanges(product types.Product) {
	if !connected {
		return
	}

	for _, user := range product.Users {
		if user.ChatType == "telegram" {

			var sizes []string
			for _, size := range product.Sizes {
				if size.Available {
					sizes = append(sizes, fmt.Sprintf("*%s*", size.Name))
				}
			}

			msgText := fmt.Sprintf(
				"[%s](%s)\nPrice: *%s*\nSizes: %s",
				product.Name,
				product.URL,
				product.Price,
				strings.Join(sizes, ", "),
			)

			chatID, _ := strconv.ParseInt(user.ChatID, 10, 64)
			message := tgbotapi.NewMessage(chatID, msgText)
			message.ParseMode = "markdown"
			bot.Send(message)
		}
	}

}

func list(msg *tgbotapi.Message) {
	products, _ := db.AllProducts()
	currentChatID := strconv.FormatInt(msg.Chat.ID, 10)

	userProducts := make([]types.Product, 0)
	for _, product := range products {
		for _, user := range product.Users {
			if user.ChatID == currentChatID {
				userProducts = append(userProducts, product)
			}
		}
	}

	if len(userProducts) == 0 {
		bot.Send(tgbotapi.NewMessage(msg.Chat.ID, "You are not tracking any products currently. Use /add to start tracking something."))
		return
	}

	msgText := "Your currently tracking products:\n"
	for _, product := range userProducts {
		msgText = msgText + fmt.Sprintf("_%d_ - [%s](%s)\n", product.ID, product.Name, product.URL)
	}

	message := tgbotapi.NewMessage(msg.Chat.ID, msgText)
	message.DisableWebPagePreview = true
	message.ParseMode = "markdown"
	bot.Send(message)
}

func rm(msg *tgbotapi.Message) {
	productID := strings.Split(msg.Text, "/rm ")[1]
	productID = strings.Trim(productID, " \n")
	err := db.DeleteProductByID(productID)
	if err != nil {
		return
	}

	bot.Send(tgbotapi.NewMessage(msg.Chat.ID, fmt.Sprintf("Product %s is not being tracked anymore.", productID)))
}

func help(msg *tgbotapi.Message) {
	if !connected {
		return
	}

	msgText := "Avalable commands:\n" +
		"/add <url> - add new product url to watch list\n" +
		"/list - shows all tracked items\n" +
		"/rm <id> - removes item with <id> (from a list) from tracking\n"

	message := tgbotapi.NewMessage(msg.Chat.ID, msgText)
	bot.Send(message)
}

func logger() *logrus.Entry {
	return logrus.WithFields(logrus.Fields{
		"module": "integrations.telegam",
	})
}
