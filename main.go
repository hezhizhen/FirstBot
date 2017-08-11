package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
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
					txt = fmt.Sprintf("格式输入错误哦。栗子：/countdown 10")
				} else {
					time.Sleep(time.Second * time.Duration(second))
					txt = fmt.Sprintf("时间到啦!")
				}
			case "time":
				txt = fmt.Sprintf("现在是：%s", rmsg.Time().In(time.Local))
			case "weather":
				txt = fmt.Sprintf("当前地区的天气信息未知")
			case "help":
				txt = fmt.Sprintf(`可用指令:
				- /start : 开始使用这个Bot
				- /countdown : 倒计时
				- /time: 现在几点了
				- /weather: 今天的天气如何
				- /help : 帮助信息`)
			default:
				txt = fmt.Sprintf("Invalid command. Please enter /help for available commands")
			}
		} else {
			txt = classifyContent(rmsg.Text)
		}

		msg := tgbotapi.NewMessage(rmsg.Chat.ID, txt)
		bot.Send(msg)
		fmt.Printf("[%s] %s\n", bot.Self.UserName, msg.Text)
	}
}

func classifyContent(s string) string {
	switch {
	case strings.Contains(s, "你好"): // greeting
		return "你好呀！"
	case strings.Contains(s, "Big Brother"): // Big Brother
		return "Big Brother is watching you!"
	case strings.Contains(s, "荤段子"): // dirty talk
		return "你是王钊吗？"
	case strings.Contains(s, "笑话"): // joke
		return "诸葛亮通过标志重捕法发现巴蜀盛产孟获"
	default:
		return "你的来信已收到"
	}
}
