package TelegramBot

import (
	"TelegramBot/Clients/Telegram"
	"flag"
	"log"
)

const (
	tgBotHost = "https://api.telegram.org"
)

func mustToken() string {
	token := flag.String("token-bot",
		"",
		"telegram bot token")

	flag.Parse()

	if *token == " " {
		log.Fatal("token is empty")
	}
	return *token
}

func main() {
	tgClient := Telegram.New(mustToken(tgBotHost))

}
