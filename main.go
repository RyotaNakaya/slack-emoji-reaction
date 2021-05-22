package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/RyotaNakaya/slack-emoji-reaction/lib"
	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("hello, world")
	// lib.GetConversations()
	// lib.GetEmoji()

	// 一旦決め打ちで
	// TODO: duration は外から指定できるようにする
	s := time.Date(2021, time.April, 01, 00, 00, 00, 0, time.UTC).Unix()
	e := time.Date(2021, time.May, 01, 00, 00, 00, 0, time.UTC).Unix()
	lib.GetChannelHistory([]string{"CF724P8RE"}, int(e), int(s))
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file, error: %v", err)
	}
	lib.SLACK_BOT_TOKEN = os.Getenv("SLACK_BOT_TOKEN")
	lib.SLACK_USER_TOKEN = os.Getenv("SLACK_USER_TOKEN")
}
