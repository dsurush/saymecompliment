package main

import (
	"context"
	"fmt"
	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"github.com/jackc/pgx/pgxpool"
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

type GetListID struct {
	ID int64
	Name string
}

func GetKillList(pool *pgxpool.Pool) (Compliments []GetListID, err error){
	conn, err := pool.Acquire(context.Background())
	if err != nil {
		log.Printf("can't get connection %e", err)
	}
	defer conn.Release()

	rows, err := conn.Query(context.Background(), `Select id, name from compliments where type = 1`)
	if err != nil {
		fmt.Printf("can't read user rows %e", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		Compliment := GetListID{}
		err := rows.Scan(
			&Compliment.ID,
			&Compliment.Name,
		)
		if err != nil {
			fmt.Println("can't scan err is = ", err)
		}
		Compliments = append(Compliments, Compliment)
	}
	if rows.Err() != nil {
		log.Printf("rows err %s", err)
		return nil, rows.Err()
	}
	return
}

func getList2(array []GetListID) string {
	list := ``
	for i := 0; i < len(array); i++{
		list = fmt.Sprintf("%s%d. %d %s\n", list, i+1, array[i].ID, array[i].Name)
		//		list+=array[i] + "\n"
	}
	return list
}

func getList(array []string) string {
	list := ``
	for i := 0; i < len(array); i++{
		list = fmt.Sprintf("%s%d. %s\n", list, i+1, array[i])
//		list+=array[i] + "\n"
	}
	return list
}
func KillByID(id string, pool *pgxpool.Pool) (err error){
	conn, err := pool.Acquire(context.Background())
	if err != nil {
		log.Printf("can't get connection %e", err)
	}
	defer conn.Release()

	_, err = conn.Exec(context.Background(), `update compliments set type = 0 where id = ($1)`, id)
	if err != nil {
		fmt.Printf(" Cant Kill %e", err)
		return
	}
	return
}

func AddComplimet(compliment string, pool *pgxpool.Pool) (err error) {
	conn, err := pool.Acquire(context.Background())
	if err != nil {
	log.Printf("can't get connection %e", err)
	}
	defer conn.Release()

	_, err = conn.Exec(context.Background(), `Insert into "compliments"(name, type) values(($1), 1)`, compliment)
	if err != nil {
		fmt.Printf(" Cant Get %e", err)
		return
	}
	return
}

func GetCompliments(pool *pgxpool.Pool) (Compliments []string, err error) {
	conn, err := pool.Acquire(context.Background())
	if err != nil {
		log.Printf("can't get connection %e", err)
	}
	defer conn.Release()

	rows, err := conn.Query(context.Background(), `Select name from compliments where type = 1`)
	if err != nil {
		fmt.Printf("can't read user rows %e", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		Compliment := ``
		err := rows.Scan(
			&Compliment,
		)
		if err != nil {
			fmt.Println("can't scan err is = ", err)
		}
		Compliments = append(Compliments, Compliment)
	}
	if rows.Err() != nil {
		log.Printf("rows err %s", err)
		return nil, rows.Err()
	}
	return
}
func SaveID(ID int64, pool *pgxpool.Pool) {
///TODO:
	conn, err := pool.Acquire(context.Background())
	if err != nil {
		log.Printf("can't get connection %e", err)
	}
	defer conn.Release()

	_, err = conn.Exec(context.Background(), `Insert into "users"(chat_id) values(($1))`, ID)
	if err != nil {
		fmt.Printf(" Cant Get %e", err)
		return
	}
	return
}
func SayNewAddCompliment(pool *pgxpool.Pool) ([]int64, error){
	conn, err := pool.Acquire(context.Background())
	if err != nil {
		log.Printf("can't get connection %e", err)
	}
	defer conn.Release()

	rows, err := conn.Query(context.Background(), `Select chat_id from users`)
	if err != nil {
		fmt.Printf("can't read user rows %e", err)
		return nil, err
	}
	defer rows.Close()
	var chats []int64
	for rows.Next() {
		var chat int64
		err := rows.Scan(
			&chat,
		)
		if err != nil {
			fmt.Println("can't scan err is = ", err)
		}
		chats = append(chats, chat)
	}
	if rows.Err() != nil {
		log.Printf("rows err %s", err)
		return nil, err
	}

	return chats, err
}
func main() {

	pool, err := pgxpool.Connect(context.Background(), `postgres://dsurush:dsurush@localhost:5432/tgtest?sslmode=disable`)
	//pool, err := pgxpool.Connect(context.Background(), `postgres://dsurush:dsurush@172.16.7.252:5432/tgtest?sslmode=disable`)
	if err != nil {
		log.Printf("Owibka - %e", err)
		log.Fatal("Can't Connection to DB")
	} else {
		fmt.Println("CONNECTION TO DB IS SUCCESS")
	}

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
//	compliments := []string{`Утонуть бы в твоих глазах.`, `От тебя несет теплотой.`, `Однажды ты спросила, что я люблю больше всего, я ответил "кофе", ты ушла, так и не услышав "с тобой"`, `Твой характер не 'американо'`, `Все твои проблемы объясняются тем, что ты просто СКОРПИОН`}

	for update := range updates {
		// универсальный ответ на любое сообщение
		compliments, err := GetCompliments(pool)
		if err != nil {
			fmt.Println("Не получается взять из бд комплименты")
		}
		list := getList(compliments)
		killLists, err := GetKillList(pool)
		if err != nil {
			fmt.Println("Не получается взять из бд комплименты c kill")
		}
		killList := getList2(killLists)

		reply := `Я не понимаю о чем ты, но ты сегодня прекрасна как никогда`
		if update.Message == nil {
			continue
		}
		// логируем от кого какое сообщение пришло
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		text := update.Message.Text
		lentext := len(text)
		if	lentext > 10 && text[0:4] == `*bb*` && text[lentext - 4: lentext] == `*bb*`{
//			compliments = append(compliments, )
			err = AddComplimet(text[4: lentext - 4], pool)
			if err != nil {
				fmt.Println("пошло не так")
				continue
				//return
			}

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, `Ваш комплимент успешно добавлен`)
			// отправляем

			chats, err := SayNewAddCompliment(pool)
			if err != nil {
				fmt.Println("can't get chats")
			} else {
				for _, value := range chats{
					message := tgbotapi.NewMessage(value, text[4:lentext-4])
					botAPI.Send(message)
				}
			}
			botAPI.Send(msg)
			continue
		}
		//fmt.Println(text[0:5], text[5:])
		if lentext > 5 && text[0:5] == `kill `{
			err := KillByID(text[5:], pool)
			if err != nil {
				fmt.Println("Не удалось удалить комплимент")
				continue
			}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, `Ваш комплимент успешно удален`)
			// отправляем
			botAPI.Send(msg)
			continue
		}
		// свитч на обработку комманд
		// комманда - сообщение, начинающееся с "/"
		switch update.Message.Command() {
		case "start":
			reply = "Привет. Я просто бот, но буду делать тебе комплименты лучше чем твой парень"
			SaveID(update.Message.Chat.ID, pool)
		case "hello":
			reply = "Hello Bitch"
		case "compliment":
			reply = compliments[rand.Intn(len(compliments))]
		case "ripoff":
			reply = ``
		case "list":
			reply = list
		case "new":
			reply = `1. Чтобы добавить новый комплимент нужно в начало и конец комплимента добавить ключевое слово *bb*. 

- Example: *bb*Чувства можно не скрывать, вдруг проснешься без них*bb*

2. Чтобы добавить новую фразу нужно в начало и конец фразы добавить ключевое слово *bc*.

- Example: *bc*Я другому отдана и всю жизнь буду ему верна.*bc*`
		case "killlist":
			reply = killList
		}

		// создаем ответное сообщение
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
		// отправляем
		fmt.Println(` I am chat id - `, update.Message.Chat.ID)

		botAPI.Send(msg)
	}
}
