package main

// commands: <prefix> <command (1/2 word)> <argument>
// commands list:
//      бабуля, какие у меня локации
//      бабуля, последние новости по дом балкон
//      бабуля, все новости по дом балкон

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	_ "github.com/mattn/go-sqlite3" // some
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/net/proxy"
)

// SystemConfig contains system config
var SystemConfig ConfigData

// ProxyAwareHTTPClient for using proxy
func ProxyAwareHTTPClient() *http.Client {
	var dialer proxy.Dialer
	dialer = proxy.Direct
	// read env and, if set proxy, apply
	proxyServer, isSet := os.LookupEnv("HTTP_PROXY")
	if isSet {
		proxyURL, err := url.Parse(proxyServer)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Invalid proxy url %q\n", proxyURL)
		}
		dialer, err = proxy.FromURL(proxyURL, proxy.Direct)
		_check(err)
	}
	// setup a http client
	httpTransport := &http.Transport{}
	httpClient := &http.Client{Transport: httpTransport}
	httpTransport.Dial = dialer.Dial
	return httpClient
}

func _check(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func main() {
	log.Printf("Start")

	SystemConfig, err := loadConfigFromEnv()
	_check(err)

	client := ProxyAwareHTTPClient()
	// create bot using token, client
	bot, err := tgbotapi.NewBotAPIWithClient(SystemConfig.botAPIKey, client)
	_check(err)

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":2112", nil)
	}()

	// debug mode on
	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	// set update interval
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 1000

	updates, err := bot.GetUpdatesChan(u)
	_check(err)

	var message InternalMessage

	// get new updates
	for update := range updates {
		// if message from channel
		if update.ChannelPost != nil {
			log.Println("Have message channel")
		}
		// if message from user
		if update.Message != nil {
			log.Println("Have message user")

			message = InternalMessage{
				update.Message.Chat.ID,
				update.Message.MessageID,
				update.Message.From.UserName,
				update.Message.Text,
			}
			log.Println("Message: " + message.messageText)
			start := time.Now()
			messageParse(&message, bot)
			elapsed := float64(time.Since(start))
			queryProcessingTime.Set(elapsed)
		}
	}
}
