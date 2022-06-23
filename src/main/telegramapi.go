package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type user struct {
	name   string
	msgu   string
	replay string
	msgID  int
}

func main() {
	users := make(map[string]user)

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

	go handle(&users, bot, &updates)

	for update := range updates {
		if update.Message == nil { // ignore agony non-Message Updates
			continue
		}

		if update.Message.ReplyToMessage != nil {
			for m, n := range users {
				if m == update.Message.ReplyToMessage.Text[1:len(m)+1] {
					n.msgID = update.Message.MessageID
					n.replay = update.Message.Text
					users[m] = n
					fmt.Println(users)
					break
				}
				fmt.Println(users)
			}
		}

		fmt.Println(update.Message.MessageID)

		//cmd
		if update.Message.IsCommand() {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
			switch update.Message.Command() {

			case "start":
				fmt.Println("start")

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
	fmt.Println(inus[2:])
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
func handle(users *map[string]user, bot *tgbotapi.BotAPI, updates *tgbotapi.UpdatesChannel) {
	fmt.Println(users)

	for update := range *updates {
		for {
			var le int
			choeruser := writea()
			text := writea()
			for k, v := range *users {
				if k == choeruser[:len(k)] {
					le = len(k)
					v.msgu = text
					(*users)[k] = v
					break
				}
			}
			txt := "/" + choeruser[:le] + " : " + text

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, txt)

			if (*users)[choeruser[0:le]].msgID == 0 {
				bot.Send(msg)

			} else {
				msg.ReplyToMessageID = (*users)[choeruser[0:le]].msgID
				bot.Send(msg)
			}

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

	user1 := user{"user1", "", "", 0}
	user2 := user{"user2", "", "", 0}
	(*users)["user1"] = user1
	(*users)["user2"] = user2
}
func deluser(users *map[string]user, user string) {
	delete(*users, user)
}
