package main

import (
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"github.com/mrDuderino/tbot/internal/app/commands"
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
	commander := commands.NewCommander(bot, productService)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		switch update.Message.Command() {
		case "help":
			commander.Help(update.Message)
		case "list":
			commander.List(update.Message)
		default:
			commander.Default(update.Message)
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
	}
}
