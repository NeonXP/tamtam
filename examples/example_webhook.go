// +build ignore

/**
 * Webhook example
 */
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/neonxp/tamtam"
	"github.com/neonxp/tamtam/schemes"
)

func main() {
	// Initialisation
	api := tamtam.New(os.Getenv("TOKEN"))
	host := os.Getenv("HOST")

	// Some methods demo:
	info, err := api.Bots.GetBot()
	log.Printf("Get me: %#v %#v", info, err)

	subs, _ := api.Subscriptions.GetSubscriptions()
	for _, s := range subs.Subscriptions {
		_, _ = api.Subscriptions.Unsubscribe(s.Url)
	}
	subscriptionResp, err := api.Subscriptions.Subscribe(host+"/webhook", []string{})
	log.Printf("Subscription: %#v %#v", subscriptionResp, err)

	ch := make(chan interface{}) // Channel with updates from TamTam

	http.HandleFunc("/webhook", api.GetHandler(ch))
	go func() {
		for {
			upd := <-ch
			log.Printf("Received: %#v", upd)
			switch upd := upd.(type) {
			case *schemes.MessageCreatedUpdate:
				err := api.Messages.Send(
					tamtam.NewMessage().
						SetUser(upd.Message.Sender.UserId).
						SetText(fmt.Sprintf("Hello, %s! Your message: %s", upd.Message.Sender.Name, upd.Message.Body.Text)),
				)
				log.Printf("Answer: %#v", err)
			default:
				log.Printf("Unknown type: %#v", upd)
			}
		}
	}()

	http.ListenAndServe(":10888", nil)
}
