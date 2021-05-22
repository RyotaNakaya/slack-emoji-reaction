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
	aggregateReaction("CF724P8RE", int(e), int(s))
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file, error: %v", err)
	}
	lib.SLACK_BOT_TOKEN = os.Getenv("SLACK_BOT_TOKEN")
	lib.SLACK_USER_TOKEN = os.Getenv("SLACK_USER_TOKEN")
}

func aggregateReaction(ChannelID string, latest int, oldest int) {
	// リアクション
	reactionDict := map[string]int{}
	// スレッドタイムスタンプ
	var ts []string

	messages := lib.FetchChannelMessages("CF724P8RE", latest, oldest)

	// reaction 集計
	for _, message := range messages {
		for _, v := range message.Reactions {
			if val, ok := reactionDict[v.Name]; ok {
				reactionDict[v.Name] = val + v.Count
			} else {
				reactionDict[v.Name] = v.Count
			}
		}

		// ThreadTimestamp がある場合スレッド取得で使うので溜めておく
		if t := message.Msg.ThreadTimestamp; t != "" {
			ts = append(ts, t)
		}
	}

	// スレッドを取得する
	messages = lib.FetchChannelThreadMessages("CF724P8RE", ts, latest, oldest)
	// reaction 集計
	for _, message := range messages {
		for _, v := range message.Reactions {
			if val, ok := reactionDict[v.Name]; ok {
				reactionDict[v.Name] = val + v.Count
			} else {
				reactionDict[v.Name] = v.Count
			}
		}
	}

	var filename string = "tmp/msg.txt"
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	for k, v := range reactionDict {
		row := fmt.Sprintf("%s: %d\n", k, v)
		_, err := file.WriteString(row)
		if err != nil {
			log.Fatal(err)
		}
	}
}
