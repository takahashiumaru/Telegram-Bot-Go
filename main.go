package main

import (
	"log"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	// Inisialisasi bot Telegram
	bot, err := tgbotapi.NewBotAPI("6293769087:AAFPr1lObI0P5JlMYY7sm35R2q0SI2PKcLk")
	if err != nil {
		log.Panic(err)
	}

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	updates, err := bot.GetUpdatesChan(updateConfig)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		userMessage := update.Message.Text
		now := time.Now()
		if strings.HasPrefix(userMessage, "/hi") || strings.HasPrefix(userMessage, "/halo") {
			reply := ""
			hour := now.Hour()
			if hour >= 6 && hour < 12 {
				reply = "Selamat pagi, " + update.Message.From.UserName + "!"
			} else if hour >= 12 && hour < 15 {
				reply = "Selamat siang, " + update.Message.From.UserName + "!"
			} else if hour >= 15 && hour < 18 {
				reply = "Selamat sore, " + update.Message.From.UserName + "!"
			} else {
				reply = "Selamat malam, " + update.Message.From.UserName + "!"
			}
			sendMessage(bot, update.Message.Chat.ID, reply)
		}
		if strings.HasPrefix(userMessage, "/id") {
			reply := "Ini adalah id telegram anda , " + strconv.FormatInt(int64(update.Message.From.ID), 10) + "!"

			sendMessage(bot, update.Message.Chat.ID, reply)
		}
		if strings.HasPrefix(userMessage, "/ceksaldo") {
			reply := "Kamu miskin tpi mau cek saldo !"
			sendMessage(bot, update.Message.Chat.ID, reply)
		}
		if strings.HasPrefix(userMessage, "/nazwa") {
			reply := "Anak kpopers yg harus di hilangkan !"
			sendMessage(bot, update.Message.Chat.ID, reply)
		}
	}
}

func sendMessage(bot *tgbotapi.BotAPI, chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	_, err := bot.Send(msg)
	if err != nil {
		log.Panic(err)
	}
}
