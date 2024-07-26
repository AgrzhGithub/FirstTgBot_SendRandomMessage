package main

import (
	TgClient "TelegramBot/clients/Telegram"
	event_consumer "TelegramBot/consumer/event-consumer"
	"TelegramBot/events/Telegram"
	"TelegramBot/storage/files"
	"flag"
	"log"
)

const (
	tgBotHost   = "api.telegram.org"
	storagePath = "files_storage"
	batchSize   = 100
)

func main() {
	eventsProcessor := Telegram.New(
		TgClient.New(tgBotHost, mustToken()),
		files.New(storagePath),
	)

	log.Print("service started")

	consumer := event_consumer.New(eventsProcessor, eventsProcessor, batchSize)
	if err := consumer.Start(); err != nil {
		log.Fatal("service is stoped", err)
	}
}

func mustToken() string {
	token := flag.String(
		"tg-bot-token",
		"",
		"telegram bot token")

	flag.Parse()

	if *token == " " {
		log.Fatal("token is empty")
	}
	return *token
}
