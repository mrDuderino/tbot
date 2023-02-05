package main

import (
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"github.com/mrDuderino/tbot/internal/service/product"
)

//const token = "6086163787:AAGRiMC6GjzAsdiXS9BIlRFA8h1dEHd4IQo"

func main() {
	godotenv.Load()
	token := os.Getenv("TOKEN")
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.UpdateConfig{
		Offset:  0,
		Timeout: 60,
	}

	updates := bot.GetUpdatesChan(u)

	productService := product.NewService()

	for update := range updates {
		if update.Message == nil {
			continue
		}

		switch update.Message.Command() {
		case "help":
			helpCommand(bot, update.Message)
		case "list":
			listCommand(bot, update.Message, productService)
		default:
			defaultBehavior(bot, update.Message)
		}
	}
}

func helpCommand(bot *tgbotapi.BotAPI, inputMessage *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(inputMessage.Chat.ID, "This is a test bot developed on GO")
	bot.Send(msg)
}

func listCommand(bot *tgbotapi.BotAPI, inputMessage *tgbotapi.Message, productService *product.Service) {
	outputMsgText := "Here all products:\n\n"
	products := productService.List()
	for _, p := range products {
		outputMsgText += p.Title + "\n"
	}
	msg := tgbotapi.NewMessage(inputMessage.Chat.ID, outputMsgText)
	bot.Send(msg)
}

func defaultBehavior(bot *tgbotapi.BotAPI, inputMessage *tgbotapi.Message) {
	log.Printf("[%s] %s", inputMessage.From.UserName, inputMessage.Text)

	msg := tgbotapi.NewMessage(inputMessage.Chat.ID, inputMessage.Text)
	msg.ReplyToMessageID = inputMessage.MessageID
	bot.Send(msg)
}
