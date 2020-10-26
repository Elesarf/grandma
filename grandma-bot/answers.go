package main

import (
	"bytes"
	"image/jpeg"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/icza/mjpeg"
	"github.com/nfnt/resize"
)

var enterString = "бабуля,"

var questionStrings = []string{
	"где мои глаза?",
	"последние новости по",
	"все новости по",
	"знакомься",
	"я поссорился с",
	"кто мои друзья?",
	"какие места у",
	"расскажи чо нового у",
	"у кого я в друзьях?",
	"help",
	"подробнее по",
}

func videoFromLocation(message *InternalMessage, bot *tgbotapi.BotAPI) {
	var answerString string

	log.Println("Process collected video")
	messageText := message.messageText
	log.Println("Message text: " + messageText)
	messageText = strings.Replace(messageText, enterString, "", 1)
	log.Println("Message without enter sting: " + messageText)
	messageText = strings.Replace(messageText, questionStrings[10], "", 1)
	log.Println("Message without command sting: " + messageText)
	messageText = strings.Replace(messageText, " ", "", 10)
	log.Println("Message behind trim: " + messageText)
	commandSplit := strings.Split(messageText, ":")
	if len(commandSplit) == 2 {
		mainLoc := commandSplit[0]
		addLoc := commandSplit[1]
		paths := getEventsListFromNameAndLocation(message.userName, InternalLocation{mainLoc, addLoc})
		// Video size: 200x100 pixels, FPS: 2
		aw, err := mjpeg.New("/tmp/test.avi", 400, 200, 2)
		if err != nil {
			answerString = "бабуля нихерашеньки не видела. Паходу ослепла."
		} else {

			start := time.Now()
			// Create a movie from images: 1.jpg, 2.jpg, ..., 10.jpg
			for i := range paths {
				data, err := os.Open(paths[len(paths)-1-i])
				if err != nil {
					answerString = "бабуля нихерашеньки не видела. Паходу ослепла."
				} else {
					img, err := jpeg.Decode(data)
					if err != nil {
						log.Fatal(err)
					}
					data.Close()
					m := resize.Resize(480, 0, img, resize.Lanczos3)
					buf := &bytes.Buffer{}
					jpeg.Encode(buf, m, nil)

					aw.AddFrame(buf.Bytes())
				}
			}
			elapsed := float64(time.Since(start) / 1000000)
			queryVideoProcessingTime.Set(elapsed)
			log.Println("Closing video file")
			err := aw.Close()
			if err != nil {
				answerString = "бабуля нихерашеньки не видела. Паходу ослепла."
			} else {
				vid, err := ioutil.ReadFile("/tmp/test.avi")
				if err != nil {
					log.Panic("Error: read video file error")
				}
				bytes := tgbotapi.FileBytes{Name: "video.avi", Bytes: vid}
				m := tgbotapi.NewVideoUpload(message.chatID, bytes)
				m.ReplyToMessageID = message.messageID
				bot.Send(m)
			}
		}
		if len(paths) > 0 {
			answerString = "Пользуйся внучек!"
		} else {
			answerString = "бабуля нихерашеньки не видела. Паходу ослепла."
		}
	} else {
		answerString = "бабуля нихерашеньки не поняла. Паходу крыша съехала."
	}

	msg := tgbotapi.NewMessage(message.chatID, answerString)

	msg.ParseMode = "HTML"
	msg.ReplyToMessageID = message.messageID

	bot.Send(msg)
}

func locationsList(message *InternalMessage, bot *tgbotapi.BotAPI) {
	log.Println("Start process location list message': ")
	var answerString string
	locations := getLocations(message.userName)
	if len(locations) <= 0 {
		answerString = "Бабуля не знает, внучек"
	}

	for mainLoc, addLoc := range locations {
		answerString = answerString + "\n" + mainLoc + " : " + addLoc
	}

	log.Println("Answer string: " + answerString + " locations len " + strconv.FormatInt(int64(len(locations)), 10))

	msg := tgbotapi.NewMessage(message.chatID, answerString)

	msg.ParseMode = "HTML"
	msg.ReplyToMessageID = message.messageID

	bot.Send(msg)
}

func lastEvent(message *InternalMessage, bot *tgbotapi.BotAPI) {

	var answerString string

	messageText := message.messageText
	log.Println("Message text: " + messageText)
	messageText = strings.Replace(messageText, enterString, "", 1)
	log.Println("Message without enter sting: " + messageText)
	messageText = strings.Replace(messageText, questionStrings[1], "", 1)
	log.Println("Message without command sting: " + messageText)
	messageText = strings.Replace(messageText, " ", "", 10)
	log.Println("Message behind trim: " + messageText)
	commandSplit := strings.Split(messageText, ":")

	if len(commandSplit) == 2 {
		mainLoc := commandSplit[0]
		addLoc := commandSplit[1]
		log.Println("Main location: " + mainLoc)
		log.Println("Add location: " + addLoc)

		path, timestamp := getNewestPath(message.userName, InternalLocation{mainLoc, addLoc})
		if len(path) != 0 {
			pic, err := ioutil.ReadFile(path)
			if err != nil {
				answerString = "бабуля нихерашеньки не видела. Паходу крыша съехала."
			} else {
				bytes := tgbotapi.FileBytes{Name: "image.jpg", Bytes: pic}
				m := tgbotapi.NewPhotoUpload(message.chatID, bytes)
				m.ReplyToMessageID = message.messageID
				bot.Send(m)
				mf := tgbotapi.NewDocumentUpload(message.chatID, bytes)
				mf.ReplyToMessageID = message.messageID
				bot.Send(mf)
				tm := time.Unix(timestamp, 0)
				answerString = "бабуля видела от такое в " + tm.String()
				photoQuery.Inc()
			}
		} else {
			answerString = "бабуля нихерашеньки не видела. Паходу ослепла."
		}
	} else {
		answerString = "бабуля нихерашеньки не поняла. Паходу крыша съехала."
	}

	msg := tgbotapi.NewMessage(message.chatID, answerString)

	msg.ParseMode = "HTML"
	msg.ReplyToMessageID = message.messageID

	bot.Send(msg)
}

func addFriend(message *InternalMessage, bot *tgbotapi.BotAPI) {
	log.Println("Start add friend")
	var answerString string
	messageText := message.messageText
	log.Println("Message text: " + messageText)
	messageText = strings.Replace(messageText, enterString, "", 1)
	log.Println("Message without enter sting: " + messageText)
	messageText = strings.Replace(messageText, questionStrings[3], "", 1)
	log.Println("Message without command sting: " + messageText)
	messageText = strings.Replace(messageText, " ", "", 10)
	log.Println("Message behind trim: " + messageText)

	if len(strings.Fields(messageText)) != 1 {
		answerString = "Бабуля думает ты бредишь"
	} else {
		err := addFriendToDb(message.userName, messageText)
		if err != nil {
			answerString = "Бабуле он не нравится"
			log.Println(err)
		} else {
			answerString = "Бабуля одобряет"
		}
	}

	msg := tgbotapi.NewMessage(message.chatID, answerString)

	msg.ParseMode = "HTML"
	msg.ReplyToMessageID = message.messageID

	bot.Send(msg)
}

func removeFriend(message *InternalMessage, bot *tgbotapi.BotAPI) {
	log.Println("Start remove friend")
	var answerString string
	messageText := message.messageText
	log.Println("Message text: " + messageText)
	messageText = strings.Replace(messageText, enterString, "", 1)
	log.Println("Message without enter sting: " + messageText)
	messageText = strings.Replace(messageText, questionStrings[4], "", 1)
	log.Println("Message without command sting: " + messageText)
	messageText = strings.Replace(messageText, " ", "", 10)
	log.Println("Message behind trim: " + messageText)

	if len(strings.Fields(messageText)) != 1 {
		answerString = "Бабуля думает ты бредишь"
	} else {
		err := removeFriendFromDb(message.userName, messageText)
		if err != nil {
			answerString = "Бабуле он нравится, не стану забывать"
			log.Println(err)
		} else {
			answerString = "Бабуля расстроена,что вы больше не дружите, но бабуля вам не указ."
		}
	}

	msg := tgbotapi.NewMessage(message.chatID, answerString)

	msg.ParseMode = "HTML"
	msg.ReplyToMessageID = message.messageID

	bot.Send(msg)
}

func friendsList(message *InternalMessage, bot *tgbotapi.BotAPI) {

	log.Println("Start remove friend")
	var answerString string

	friends := getFriendList(message.userName)

	if len(friends) > 0 {
		answerString = "бабуля знает, что ты дружишь с:\n"
		for _, friend := range friends {
			answerString = answerString + friend + "\n"
		}
	} else {
		answerString = "ты паходу интроверт, дружок"
	}

	msg := tgbotapi.NewMessage(message.chatID, answerString)

	msg.ParseMode = "HTML"
	msg.ReplyToMessageID = message.messageID

	bot.Send(msg)
}

func locationsOfFriend(message *InternalMessage, bot *tgbotapi.BotAPI) {
	log.Println("Start locations of friend")
	var answerString string
	messageText := message.messageText
	log.Println("Message text: " + messageText)
	messageText = strings.Replace(messageText, enterString, "", 1)
	log.Println("Message without enter sting: " + messageText)
	messageText = strings.Replace(messageText, questionStrings[6], "", 1)
	log.Println("Message without command sting: " + messageText)
	messageText = strings.Replace(messageText, " ", "", 10)
	log.Println("Message behind trim: " + messageText)

	if thatUserIsFriend(message.userName, messageText) {
		locations := getLocations(messageText)
		answerString = "бабуля считает что " + message.userName + " тебя знает.\n"
		answerString = answerString + "бабуля скажет:"

		for location, addLoc := range locations {
			answerString = answerString + "\n" + location + " : " + addLoc
		}
	} else {
		answerString = "бабуля думает, что вы не дружите. Бабуля помолчит."
	}

	msg := tgbotapi.NewMessage(message.chatID, answerString)
	msg.ParseMode = "HTML"
	msg.ReplyToMessageID = message.messageID

	bot.Send(msg)
}

func lastEventOfFriend(message *InternalMessage, bot *tgbotapi.BotAPI) {
	log.Println("Start last event of friend")
	// var answerString string
	var answerString string
	messageText := message.messageText
	log.Println("Message text: " + messageText)
	messageText = strings.Replace(messageText, enterString, "", 1)
	log.Println("Message without enter sting: " + messageText)
	messageText = strings.Replace(messageText, questionStrings[7], "", 1)
	log.Println("Message without command sting: " + messageText)
	messageText = strings.Replace(messageText, " ", "", 10)
	log.Println("Message behind trim: " + messageText)
	commandSplit := strings.Split(messageText, "по")
	if len(commandSplit) == 2 {
		log.Println("Command: " + commandSplit[0])
		log.Println("Location: " + commandSplit[1])
		userName := commandSplit[0]
		if thatUserIsFriend(message.userName, userName) {
			fakeMessage := message
			fakeMessage.userName = userName
			fakeMessage.messageText = "бабуля, последние новости по " + commandSplit[1]
			lastEvent(fakeMessage, bot)
			answerString = "бабуля слышала что у " + userName + " :"
		} else {
			answerString = "бабуле кажется, что ты не о том спрашиваешь"
		}
	} else {
		answerString = "бабуле кажется, что ты не о том спрашиваешь"
	}

	msg := tgbotapi.NewMessage(message.chatID, answerString)

	msg.ParseMode = "HTML"
	msg.ReplyToMessageID = message.messageID

	bot.Send(msg)
}

func iInFriendsOn(message *InternalMessage, bot *tgbotapi.BotAPI) {
	log.Println("Start i in friends on")
	var answerString string

	friends := getIInFriendsOn(message.userName)

	if len(friends) > 0 {
		answerString = "бабуля знает, что с тобой дружат:\n"
		for _, friend := range friends {
			answerString = answerString + friend + "\n"
		}
	} else {
		answerString = "ты паходу интроверт, дружок"
	}

	msg := tgbotapi.NewMessage(message.chatID, answerString)

	msg.ParseMode = "HTML"
	msg.ReplyToMessageID = message.messageID

	bot.Send(msg)
}

func help(message *InternalMessage, bot *tgbotapi.BotAPI) {

	answerString := "бабуля умеет: \n"

	for _, command := range questionStrings {
		answerString = answerString + command + "\n"
	}

	msg := tgbotapi.NewMessage(message.chatID, answerString)

	msg.ParseMode = "HTML"
	msg.ReplyToMessageID = message.messageID

	bot.Send(msg)
}
