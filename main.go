package main

import (
	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"log"
	"math/rand"
	"time"
)

var (
	// глобальная переменная в которой храним токен
	telegramBotToken string
)

//
//func init() {
//  // принимаем на входе флаг -telegrambottoken
//  flag.StringVar(&telegramBotToken, "telegrambottoken", "", "Telegram Bot Token")
//  flag.Parse()
//  // без него не запускаемся
//  if telegramBotToken == "" {
//    log.Print("-telegrambottoken is required")
//    os.Exit(1)
//  }
//}

func main() {
	botAPI, err := tgbotapi.NewBotAPI(`1420536268:AAGXIHTSyI7jq0PEI4oDsgN1_NFfUH9Bzcw`)
	if err != nil {
		log.Print(err)
		time.Sleep(time.Second * 10)
		return
	}
	log.Printf("Authorized on account %s", botAPI.Self.UserName)
	// u - структура с конфигом для получения апдейтов
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	// используя конфиг u создаем канал в который будут прилетать новые сообщения
	updates, err := botAPI.GetUpdatesChan(u)

	for update := range updates {
		// универсальный ответ на любое сообщение
		compliments := []string{`Утонуть бы в твоих глазах`, `От тебя несет теплотой`, `Однажды ты спросила, что я люблю больше всего, я ответил "кофе", ты ушла, так и не услышав "с тобой""`, `Твой характер не 'американо'`, `Все твои проблемы объясняется тем, что ты просто СКОРПИОН`}
		reply := compliments[rand.Intn(len(compliments))]
		if update.Message == nil {
			continue
		}

		// логируем от кого какое сообщение пришло
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		// свитч на обработку комманд
		// комманда - сообщение, начинающееся с "/"
		switch update.Message.Command() {
		case "start":
			reply = "Привет. Я просто бот, но буду делать тебе комплименты лучше чем твой парень"
		case "hello":
			reply = "Hello Bitch"
		}

		// создаем ответное сообщение
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
		// отправляем
		botAPI.Send(msg)
	}
}
