package main

import (
	"fmt"
	"log"
	"os"

	"github.com/RyotaNakaya/slack-emoji-reaction/lib"
	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("hello, world")
	// lib.GetConversations()
	// lib.GetEmoji()
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file, error: %v", err)
	}
	lib.SLACK_TOKEN = os.Getenv("SLACK_BOT_TOKEN")
}
