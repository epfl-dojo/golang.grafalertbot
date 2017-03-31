package main

import (
	"github.com/deckarep/golang-set"
	"gopkg.in/telegram-bot-api.v4"
	"log"
	"net/http"
	"os"
)

var bot, err = tgbotapi.NewBotAPI(os.Args[1])
var pageChan = make(chan string)

func main() {
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)
	registerChan := make(chan int64)
	go register(registerChan)
	go greetNewUsers(registerChan, pageChan)

	http.HandleFunc("/", handle)
	http.ListenAndServe("0.0.0.0:8000", nil)
}

func handle(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	pageChan <- "ALERT"
}

func register(registerChan chan int64) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Panic(err)
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		msg.ReplyToMessageID = update.Message.MessageID

		bot.Send(msg)
		registerChan <- update.Message.Chat.ID
	}
}

func greetNewUsers(registerChannel chan int64, pageChannel chan string) {
	lesDudes := mapset.NewSet()

	for {
		select {
		case chatID := <-registerChannel:
			log.Printf("Y'a le dude %d qui me coze", chatID)
			lesDudes.Add(chatID)

		case page := <-pageChannel:
			log.Printf("Houla, y'a le feu")
			for chatIDInterface := range lesDudes.Iter() {
				chatID, _ := chatIDInterface.(int64)
				msg := tgbotapi.NewMessage(chatID, page)
				bot.Send(msg)
			}
		}
	}
}
