package main

import (
	"fmt"
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

func getList(array []string) string {
	list := ``
	for i := 0; i < len(array); i++{
		list = fmt.Sprintf("%s%d. %s\n", list, i+1, array[i])
//		list+=array[i] + "\n"
	}
	return list
}

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
	compliments := []string{`Утонуть бы в твоих глазах.`, `От тебя несет теплотой.`, `Однажды ты спросила, что я люблю больше всего, я ответил "кофе", ты ушла, так и не услышав "с тобой"`, `Твой характер не 'американо'`, `Все твои проблемы объясняются тем, что ты просто СКОРПИОН`}

	for update := range updates {
		// универсальный ответ на любое сообщение
		list := getList(compliments)
		reply := `Я не понимаю о чем ты, но ты сегодня прекрасна как никогда`
		if update.Message == nil {
			continue
		}

		// логируем от кого какое сообщение пришло
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		text := update.Message.Text
		lentext := len(text)
		if	lentext > 10 && text[0:4] == `*bb*` && text[lentext - 4: lentext] == `*bb*`{
			compliments = append(compliments, text[4: lentext - 4])
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, `Ваш комплимент успешно добавлен`)
			// отправляем
			botAPI.Send(msg)
			continue
		}

		// свитч на обработку комманд
		// комманда - сообщение, начинающееся с "/"
		switch update.Message.Command() {
		case "start":
			reply = "Привет. Я просто бот, но буду делать тебе комплименты лучше чем твой парень"
		case "hello":
			reply = "Hello Bitch"
		case "compliment":
			reply = compliments[rand.Intn(len(compliments))]
		case "ripoff":
			reply = ``
		case "list":
			reply = list
		}

		// создаем ответное сообщение
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
		// отправляем
		botAPI.Send(msg)
	}
}
