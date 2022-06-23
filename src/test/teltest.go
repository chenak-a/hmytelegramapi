package main

import (
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/crisp-im/go-crisp-api/crisp"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type user struct {
	name    string
	msgu    string
	replay  string
	msgID   int
	chatID  int64
	session string
	time    int64
}
type usercrisp struct {
	msg       string
	msgid     int
	session   string
	websiteid string
}

func main() {
	wg := sync.WaitGroup{}
	constante := make(map[string]int64, 1)
	userscrisp := make(map[string]usercrisp, 1)

	users := make(map[string]user, 50)
	client := crisp.New()
	// set key to crisp
	client.Authenticate("9488974e-d75e-4482-a547-3807e89fcd29", "e965866eaec642658bd106954b0c5e625b406f9aa25c78caeff8f05b4156ba65")

	// set key telegram
	bot, err := tgbotapi.NewBotAPI("1168994847:AAH0AwQJAZmetuDIoKCOSHc2hR8gwrzIB8U")
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true

	//log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	//recice and send /telegram

	go handle(&users, bot, &updates, client, &userscrisp, &constante, &wg)

	for update := range updates {
		if update.Message == nil { // ignore agony non-Message Updates
			continue
		}

		if update.Message.ReplyToMessage != nil {
			// send message from telegram to website
			for m, n := range users {
				if len(update.Message.ReplyToMessage.Text) >= 5 {
					if m == update.Message.ReplyToMessage.Text[1:len(m)+1] {

						fmt.Println(m)
						fmt.Println(len(m))
						n.msgID = update.Message.MessageID
						n.replay = update.Message.Text
						users[m] = n
						client.Website.SendTextMessageInConversation(userscrisp["usercrisp1"].websiteid, n.session, crisp.ConversationTextMessageNew{Type: "text", Content: n.replay, From: "operator", Origin: "chat"})
						fmt.Println(users)
						break
					} else {
						continue
					}

				}
			}
		}
		go delusers(&users, &userscrisp, client)
		fmt.Println(update.Message.MessageID)

		//cmd
		if update.Message.IsCommand() {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
			switch update.Message.Command() {

			case "start":
				fmt.Println("start")
				constante["chatID"] = update.Message.Chat.ID
				fmt.Println(constante)
			case "menu":
				if mapuser(&users) == "" {
					msg.Text = "no user"
				} else {
					msg.Text = mapuser(&users)
				}

			case identification(users, update.Message.Text):
				usera := identification(users, update.Message.Text)
				msg.Text = "/" + usera + " : " + users[usera].msgu

			default:
				msg.Text = "I don't know that command"
			}
			bot.Send(msg)
		}

	}

}
func mapuser(users *map[string]user) string {
	var a string
	b := 0
	for k, _ := range *users {
		b++
		if b != len(*users) {
			a += "/" + k + ", "
		} else {
			a += "/" + k
		}
	}
	return a
}
func identification(users map[string]user, inus string) string {

	var fin string
	for k, _ := range users {

		if k == inus[1:len(k)+1] {
			fin = k
			break
		}
	}
	return fin
}

// handle website users creat del and send msg
func handle(users *map[string]user, bot *tgbotapi.BotAPI, updates *tgbotapi.UpdatesChannel, client *crisp.Client, userscrisp *map[string]usercrisp, constante *map[string]int64, wg *sync.WaitGroup) {
	fmt.Println("start")
	for {
		wg.Add(1)
		client.Events.Listen(
			[]string{
				"message:send",
			},

			func(reg *crisp.EventsRegister) {
				reg.On("message:send/text", func(evt crisp.EventsReceiveTextMessage) {

					// Handle text message from visitor

					msg := evt.Content
					msgid := evt.Fingerprint
					session := evt.SessionID
					websiteid := evt.EventsWebsiteGeneric.WebsiteID
					usercrisp1 := usercrisp{*msg, *msgid, *session, *websiteid}

					defer cleanup()
					go delusers(users, userscrisp, client)
					go choseinpute(&usercrisp1, userscrisp, client, updates, users, bot, constante)

				})

			},

			func() {
				// Socket is disconnected: will try to reconnect
			},

			func() {
				// Socket error: may be broken
			},
		)
		wg.Wait()
	}

}

func choseinpute(usercrisp *usercrisp, userscrisp *map[string]usercrisp, client *crisp.Client, updates *tgbotapi.UpdatesChannel, users *map[string]user, bot *tgbotapi.BotAPI, constante *map[string]int64) {

	defer cleanup()
	(*userscrisp)["usercrisp1"] = *usercrisp
	a := 0
	var b string
	if a == len(*users) {
		usern := user{"user" + strconv.Itoa(len(*users)+1), (*userscrisp)["usercrisp1"].msg, "", 0, 0, (*userscrisp)["usercrisp1"].session, time.Now().Unix()}
		(*users)["user"+strconv.Itoa(len(*users)+1)] = usern
		b = "user" + strconv.Itoa(len(*users))
	} else {
		for s, v := range *users {
			a++
			if v.session == (*userscrisp)["usercrisp1"].session {
				v.msgu = (*userscrisp)["usercrisp1"].msg
				v.time = time.Now().Unix()
				(*users)[s] = v
				b = s
				break
			} else if a == len(*users) {
				usern := user{"user" + strconv.Itoa(len(*users)+1), (*userscrisp)["usercrisp1"].msg, "", 0, 0, (*userscrisp)["usercrisp1"].session, time.Now().Unix()}
				(*users)["user"+strconv.Itoa(len(*users)+1)] = usern
				b = "user" + strconv.Itoa(len(*users))
				break
			}

		}
	}

	defer cleanup()
	// creat a user if i have to
	//handle trafic

	go sendtext(b, bot, constante, &*users)

}

func sendtext(username string, bot *tgbotapi.BotAPI, constante *map[string]int64, users *map[string]user) {

	defer cleanup()

	choeruser := (*users)[username].name

	text := (*users)[username].msgu

	txt := "/" + choeruser + " : " + text

	msg := tgbotapi.NewMessage((*constante)["chatID"], txt)

	if (*users)[choeruser].msgID == 0 {
		bot.Send(msg)

	} else {
		msg.ReplyToMessageID = (*users)[choeruser].msgID
		bot.Send(msg)

	}
}
func delusers(users *map[string]user, userscrisp *map[string]usercrisp, client *crisp.Client) {
	for s, v := range *users {
		uptime := time.Now().Unix() - v.time
		fmt.Println(uptime)
		if uptime >= 1200 {
			delete(*users, s)
		}

	}
}
func cleanup() {
	if r := recover(); r != nil {
		fmt.Println("Recovered", r)
	}
}
