package main

import (
	"fmt"
	"sync"

	"github.com/crisp-im/go-crisp-api/crisp"
)

type usercrisp struct {
	msg       string
	msgid     int
	session   string
	websiteid string
}

func main() {
	wg := sync.WaitGroup{}
	userscrisp := make(map[string]usercrisp)
	usercrisp1 := usercrisp{"", 0, "", ""}
	userscrisp["user1"] = usercrisp1
	client := crisp.New()
	// Set authentication parameters
	client.Authenticate("9488974e-d75e-4482-a547-3807e89fcd29", "e965866eaec642658bd106954b0c5e625b406f9aa25c78caeff8f05b4156ba65")

	// Connect to realtime events backend and listen (only to 'message:send' namespace)
	for {
		fmt.Println("start")

		wg.Add(1)
		client.Events.Listen(
			[]string{
				"message:send", "message:updated", "message:received",
			},

			func(reg *crisp.EventsRegister) {
				fmt.Println("connected")
				reg.On("message:send/text", func(evt crisp.EventsReceiveTextMessage) {

					// Handle text message from visitor

					msg := evt.Content
					usercrisp1.msg = *msg
					msgid := evt.Fingerprint
					usercrisp1.msgid = *msgid
					session := evt.SessionID
					usercrisp1.session = *session
					websiteid := evt.EventsWebsiteGeneric.WebsiteID
					usercrisp1.websiteid = *websiteid
					fmt.Println(*msg)
					go choseinpute(&usercrisp1, &userscrisp, client)

				})
			},

			func() {
				// Socket is disconnected: will try to reconnect
				fmt.Println("bye1")
			},

			func() {
				// Socket error: may be broken
				fmt.Println("bye2")
			},
		)
		wg.Wait()
	}
}
func choseinpute(usercrisp *usercrisp, userscrisp *map[string]usercrisp, client *crisp.Client) {

	if usercrisp.websiteid == "ff31c7cf-53f5-42b4-a87e-ea5f0b92dd15" {
		if (*userscrisp)["user1"].msg != usercrisp.msg {
			(*userscrisp)["user1"] = *usercrisp
			fmt.Println(usercrisp)
			client.Website.SendTextMessageInConversation(usercrisp.websiteid, usercrisp.session, crisp.ConversationTextMessageNew{Type: "text", Content: "I'm a bot", From: "operator", Origin: "chat"})

		}

	}
}
