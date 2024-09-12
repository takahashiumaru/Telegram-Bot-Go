package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	times "time"

	"github.com/go-resty/resty/v2"
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
	if err != nil {
		panic(err)
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}
		fmt.Println("From : " + update.Message.Chat.FirstName)
		fmt.Println("Message : " + update.Message.Text)
		userMessage := update.Message.Text
		now := times.Now()
		if strings.Contains(userMessage, "/hi") || strings.Contains(userMessage, "/halo") {
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
		} else if strings.Contains(userMessage, "/id") {
			reply := "Ini adalah id telegram anda , " + strconv.FormatInt(int64(update.Message.From.ID), 10) + "!"

			sendMessage(bot, update.Message.Chat.ID, reply)
		} else if strings.Contains(userMessage, "/ceksaldo") {
			reply := "Kamu miskin tpi mau cek saldo !"
			sendMessage(bot, update.Message.Chat.ID, reply)
		} else if strings.Contains(userMessage, "/gempa") {
			url := "https://data.bmkg.go.id/DataMKG/TEWS/autogempa.json"
			resp, err := http.Get(url)
			if err != nil {
				log.Println(err)
			}

			var data struct {
				Infogempa struct {
					Gempa struct {
						Tanggal     string `json:"Tanggal"`
						Jam         string `json:"Jam"`
						DateTime    string `json:"DateTime"`
						Coordinates string `json:"Coordinates"`
						Lintang     string `json:"Lintang"`
						Bujur       string `json:"Bujur"`
						Magnitude   string `json:"Magnitude"`
						Kedalaman   string `json:"Kedalaman"`
						Wilayah     string `json:"Wilayah"`
						Potensi     string `json:"Potensi"`
						Dirasakan   string `json:"Dirasakan"`
						Shakemap    string `json:"Shakemap"`
					} `json:"gempa"`
				} `json:"Infogempa"`
			}

			err = json.NewDecoder(resp.Body).Decode(&data)
			if err != nil {
				log.Println(err)
			}

			reply := fmt.Sprintf("Terjadi gempa pada %s jam %s dengan magnitude %s dan kedalaman %s di %s. Wilayah yang terasa adalah %s dengan skala %s. Shakemap: https://data.bmkg.go.id/DataMKG/TEWS/%s",
				data.Infogempa.Gempa.Tanggal, data.Infogempa.Gempa.Jam, data.Infogempa.Gempa.Magnitude, data.Infogempa.Gempa.Kedalaman, data.Infogempa.Gempa.Wilayah, data.Infogempa.Gempa.Dirasakan, data.Infogempa.Gempa.Potensi, data.Infogempa.Gempa.Shakemap)
			sendMessage(bot, update.Message.Chat.ID, reply)
		} else if strings.Contains(userMessage, "/listbank") {
			url := "https://api-rekening.lfourr.com/listBank"
			resp, err := http.Get(url)
			if err != nil {
				log.Println(err)
			}

			var data struct {
				Status bool   `json:"status"`
				Msg    string `json:"msg"`
				Data   []struct {
					KodeBank string `json:"kodeBank"`
					NamaBank string `json:"namaBank"`
				} `json:"data"`
			}

			err = json.NewDecoder(resp.Body).Decode(&data)
			if err != nil {
				log.Println(err)
			}

			reply := ""
			for _, bank := range data.Data {
				reply += fmt.Sprintf("%s - %s\n", bank.KodeBank, bank.NamaBank)
			}
			sendMessage(bot, update.Message.Chat.ID, reply)
		} else if strings.Contains(userMessage, "/listewalet") {
			url := "https://api-rekening.lfourr.com/listEwallet"
			resp, err := http.Get(url)
			if err != nil {
				log.Println(err)
			}

			var data struct {
				Status bool   `json:"status"`
				Msg    string `json:"msg"`
				Data   []struct {
					KodeBank string `json:"kodeBank"`
					NamaBank string `json:"namaBank"`
				} `json:"data"`
			}

			err = json.NewDecoder(resp.Body).Decode(&data)
			if err != nil {
				log.Println(err)
			}

			reply := ""
			for _, ewallet := range data.Data {
				reply += fmt.Sprintf("%s - %s\n", ewallet.KodeBank, ewallet.NamaBank)
			}
			sendMessage(bot, update.Message.Chat.ID, reply)
		} else if strings.Contains(userMessage, "/play") {
			song := strings.TrimSpace(strings.Replace(userMessage, "/play", "", 1))
			if song != "" {
				url := fmt.Sprintf("https://api.deezer.com/search?q=%s&limit=1", url.QueryEscape(song))
				resp, err := http.Get(url)
				if err != nil {
					log.Println(err)
				}

				var data struct {
					Data []struct {
						Title   string `json:"title"`
						Preview string `json:"preview"`
					} `json:"data"`
				}

				err = json.NewDecoder(resp.Body).Decode(&data)
				if err != nil {
					log.Println(err)
				}

				if len(data.Data) > 0 {
					audio := tgbotapi.NewAudioShare(update.Message.Chat.ID, data.Data[0].Preview)
					audio.Duration = 100000000 // 100000000 = 100 seconds
					_, err = bot.Send(audio)
					if err != nil {
						log.Println(err)
					}
				}
			}
		} else if strings.Contains(userMessage, "/jokesorangluar") {
			resp, err := http.Get("https://official-joke-api.appspot.com/random_joke")
			if err != nil {
				log.Println(err)
			}

			var data struct {
				Setup     string `json:"setup"`
				Punchline string `json:"punchline"`
			}

			err = json.NewDecoder(resp.Body).Decode(&data)
			if err != nil {
				log.Println(err)
			}

			setup, err := translateToIndonesian(data.Setup)
			if err != nil {
				log.Println(err)
			}

			punchline, err := translateToIndonesian(data.Punchline)
			if err != nil {
				log.Println(err)
			}

			reply := setup
			sendMessage(bot, update.Message.Chat.ID, reply)
			times.Sleep(2 * times.Second)
			reply = punchline
			sendMessage(bot, update.Message.Chat.ID, reply)
		} else if strings.Contains(userMessage, "/cuaca") {
			city := strings.TrimSpace(strings.ReplaceAll(userMessage, "/cuaca", ""))
			if city == "" {
				city = "dki-jakarta"
			}
			url := fmt.Sprintf("https://cuaca-gempa-rest-api.vercel.app/weather/%s", city)
			resp, err := http.Get(url)
			if err != nil {
				log.Println(err)
			}
			var data struct {
				Success bool   `json:"success"`
				Message string `json:"message"`
				Data    struct {
					Issue struct {
						Timestamp string `json:"timestamp"`
						Year      string `json:"year"`
						Month     string `json:"month"`
						Day       string `json:"day"`
						Hour      string `json:"hour"`
						Minute    string `json:"minute"`
						Second    string `json:"second"`
					} `json:"issue"`
					Areas []struct {
						ID         string `json:"id"`
						Latitude   string `json:"latitude"`
						Longitude  string `json:"longitude"`
						Coordinate string `json:"coordinate"`
						Type       string `json:"type"`
						Region     string `json:"region"`
						Level      string `json:"level"`
						Desc       string `json:"description"`
						Domain     string `json:"domain"`
						Tags       string `json:"tags"`
						Params     []struct {
							ID    string `json:"id"`
							Desc  string `json:"description"`
							Type  string `json:"type"`
							Times []struct {
								Type       string `json:"type"`
								H          string `json:"h"`
								Datetime   string `json:"datetime"`
								Celcius    string `json:"celcius"`
								Fahrenheit string `json:"fahrenheit"`
							} `json:"times"`
						} `json:"params"`
					} `json:"areas"`
				} `json:"data"`
			}

			err = json.NewDecoder(resp.Body).Decode(&data)
			if err != nil {
				log.Println(err)
			}

			if data.Success {
				for _, area := range data.Data.Areas {
					for _, param := range area.Params {
						if param.ID == "t" {
							for _, time := range param.Times {
								if time.H == "6" {
									day := times.Now().Format("02 January 2006")
									reply := fmt.Sprintf("Cuaca di %s pada tanggal %s : %s Celcius", area.Desc, day, time.Celcius)
									sendMessage(bot, update.Message.Chat.ID, reply)
								}
							}
						}
					}
				}
			} else {
				if strings.Contains(userMessage, "/") {
					photoBytes, err := os.ReadFile("mrsbrewc.jpg")
					if err != nil {
						log.Println(err)
					}

					photo := tgbotapi.FileBytes{
						Name:  "mrsbrewc.jpg",
						Bytes: photoBytes,
					}

					photoMessage := tgbotapi.NewPhotoUpload(update.Message.Chat.ID, photo)
					_, err = bot.Send(photoMessage)
					if err != nil {
						log.
							Println(err)
					}
					sendMessage(bot, update.Message.Chat.ID, "kamu cabul ya")
				}
			}
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

type TranslationResponse struct {
	TranslatedText string `json:"translatedText"`
}

// Fungsi untuk menerjemahkan teks dari bahasa Inggris ke bahasa Indonesia
func translateToIndonesian(text string) (string, error) {
	client := resty.New()
	apiURL := "https://api.mymemory.translated.net/get?" // URL MyMemory API

	// Membuat request query
	reqQuery := url.Values{
		"q":        {text},
		"langpair": {"en|id"},
	}

	resp, err := client.R().
		SetQueryParamsFromValues(reqQuery).
		Get(apiURL)

	if err != nil {
		return "", err
	}

	var translation struct {
		ResponseData struct {
			TranslatedText string `json:"translatedText"`
		} `json:"responseData"`
	}

	err = json.Unmarshal(resp.Body(), &translation)
	if err != nil {
		return "", err
	}

	return translation.ResponseData.TranslatedText, nil
}
