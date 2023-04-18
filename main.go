package main

import (
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	bot, err := tgbotapi.NewBotAPI("-")
	if err != nil {
		log.Panic(err)
	}

	// enable debug mode
	bot.Debug = true

	// set timeout koneksi
	bot.Client.Timeout = 30 * time.Second

	// konfigurasi update
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	// dapatkan channel untuk menerima update
	updates, err := bot.GetUpdatesChan(updateConfig)
	if err != nil {
		log.Panic(err)
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}

		// proses pesan yang diterima
		switch update.Message.Text {
		case "/start":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Halo, selamat datang di bot saya!")
			bot.Send(msg)

		case "Halo":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Halo juga!")
			bot.Send(msg)

		default:
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Maaf, saya tidak mengerti maksud Anda.")
			bot.Send(msg)
		}
	}
}
