package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	args := os.Args
	token := args[1]
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		fmt.Println(err)
	}
	// getMe
	user, err := bot.GetMe()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("bot", user.UserName, "is watching...")
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		fmt.Println(err)
	}
	for update := range updates {
		if update.Message == nil {
			continue
		}
		rmsg := update.Message
		fmt.Printf("[%s] %s\n", rmsg.From.UserName, rmsg.Text)

		var txt string
		if rmsg.IsCommand() {
			switch rmsg.Command() {
			case "start":
				txt = fmt.Sprintf("Welcome, %s", rmsg.From.UserName)
			case "countdown":
				secondPattern := regexp.MustCompile(`[0-9]+`)
				secondS := secondPattern.FindString(rmsg.Text)
				second, err := strconv.Atoi(secondS)
				if err != nil {
					fmt.Println(err)
					txt = fmt.Sprintf("fail to extract the interval")
				} else {
					time.Sleep(time.Second * time.Duration(second))
					txt = fmt.Sprintf("Time is up!")
				}
			case "help":
				txt = fmt.Sprintf(`Available commands:
				- /start : start to use this bot
				- /countdown : enter the interval(unit:s) you would like to wait
				- /help : more info of usage about this bot`)
			default:
				txt = fmt.Sprintf("Invalid command. Please enter /help for available commands")
			}
		} else {
			txt = "Your message has been received."
		}

		msg := tgbotapi.NewMessage(rmsg.Chat.ID, txt)
		bot.Send(msg)
		fmt.Printf("[%s] %s\n", bot.Self.UserName, msg.Text)
	}
}
