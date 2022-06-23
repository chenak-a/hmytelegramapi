package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/crisp-im/go-crisp-api/crisp"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type user struct {
	name   string
	msgu   string
	replay string
	msgID  int
	chatID int64
}
type usercrisp struct {
	msg       string
	msgid     int
	session   string
	websiteid string
}

func main() {
	constante := make(map[string]int64)
	userscrisp := make(map[string]usercrisp)

	users := make(map[string]user)
	client := crisp.New()
	// set key to crisp
	client.Authenticate("9488974e-d75e-4482-a547-3807e89fcd29", "e965866eaec642658bd106954b0c5e625b406f9aa25c78caeff8f05b4156ba65")
	websiteData := crisp.WebsiteCreate{
		Name:   "test",
		Domain: "onebadger.crypto",
	}
	client.Website.CreateWebsite(websiteData)
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
	// creat users

	go creatuser(&users)
	//recice and send /telegram

	go handle(&users, bot, &updates, client, &userscrisp, &constante)

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
						client.Website.SendTextMessageInConversation(userscrisp["usercrisp1"].websiteid, userscrisp["usercrisp1"].session, crisp.ConversationTextMessageNew{Type: "text", Content: n.replay, From: "operator", Origin: "chat"})
						fmt.Println(users)
						break
					} else {
						continue
					}

				}
			}
		}

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
				if mapuser(users) == "" {
					msg.Text = "no user"
				} else {
					msg.Text = mapuser(users)
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
func mapuser(users map[string]user) string {
	var a string
	b := 0
	for k, _ := range users {
		b++
		if b != len(users) {
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

/*func updateuserinfo(users *map[string]user,name string ,m user) {
	name.m =
}*/

// handle website users creat del and send msg
func handle(users *map[string]user, bot *tgbotapi.BotAPI, updates *tgbotapi.UpdatesChannel, client *crisp.Client, userscrisp *map[string]usercrisp, constante *map[string]int64) {
	fmt.Println("start")
	for {
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
	}

}

func choseinpute(usercrisp *usercrisp, userscrisp *map[string]usercrisp, client *crisp.Client, updates *tgbotapi.UpdatesChannel, users *map[string]user, bot *tgbotapi.BotAPI, constante *map[string]int64) {

	if usercrisp.websiteid == "ff31c7cf-53f5-42b4-a87e-ea5f0b92dd15" {

		if (*userscrisp)["usercrisp1"].msgid != usercrisp.msgid {

			(*userscrisp)["usercrisp1"] = *usercrisp
			fmt.Println(usercrisp)
			mmsg := usercrisp.msg
			// creat a user if i have to
			//handle trafic
			go sendtext(mmsg, bot, constante, users)

		}
	}

}

func writea() string {
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	fmt.Println(text)
	return text
}

// creat a user
func creatuser(users *map[string]user) {

	user1 := user{"user1", "", "", 0, 0}
	user2 := user{"user2", "", "", 0, 0}
	(*users)["user1"] = user1
	(*users)["user2"] = user2
}
func deluser(users *map[string]user, user string) {
	delete(*users, user)
}
func sendtext(message string, bot *tgbotapi.BotAPI, constante *map[string]int64, users *map[string]user) {
	var le int
	choeruser := "user1"
	text := message
	for k, v := range *users {
		if k == choeruser[:len(k)] {
			le = len(k)
			v.msgu = text
			(*users)[k] = v
			break
		}
	}
	txt := "/" + choeruser[:le] + " : " + text

	msg := tgbotapi.NewMessage((*constante)["chatID"], txt)

	if (*users)[choeruser[0:le]].msgID == 0 {
		bot.Send(msg)

	} else {
		msg.ReplyToMessageID = (*users)[choeruser[0:le]].msgID
		bot.Send(msg)

	}
}
