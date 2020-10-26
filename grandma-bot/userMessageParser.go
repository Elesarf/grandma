package main

import (
	"log"
	"strings"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func messageParse(message *InternalMessage, bot *tgbotapi.BotAPI) {

	log.Println("Start parse message")
	if message.messageText == "" {
		log.Println("Error parse: message empty")
		return
	}

	message.messageText = strings.ToLower(message.messageText)
	message.messageText = strings.TrimRight(message.messageText, " ")
	message.messageText = strings.TrimLeft(message.messageText, " ")

	if !strings.Contains(message.messageText, enterString) {
		log.Println("Error parse: message not contain enter string")
		return
	}

	queryCount.Inc()

	switch {
	case strings.Contains(message.messageText, questionStrings[0]):
		go locationsList(message, bot)
		return
	case strings.Contains(message.messageText, questionStrings[1]):
		go lastEvent(message, bot)
		return
	case strings.Contains(message.messageText, questionStrings[3]):
		addFriend(message, bot)
		return
	case strings.Contains(message.messageText, questionStrings[4]):
		go removeFriend(message, bot)
		return
	case strings.Contains(message.messageText, questionStrings[5]):
		go friendsList(message, bot)
		return
	case strings.Contains(message.messageText, questionStrings[6]):
		go locationsOfFriend(message, bot)
		return
	case strings.Contains(message.messageText, questionStrings[7]):
		go lastEventOfFriend(message, bot)
		return
	case strings.Contains(message.messageText, questionStrings[8]):
		go iInFriendsOn(message, bot)
		return
	case strings.Contains(message.messageText, questionStrings[9]):
		go help(message, bot)
		return
	case strings.Contains(message.messageText, questionStrings[10]):
		go videoFromLocation(message, bot)
		return
	}

	log.Println("Error parse: answer not found")
}
