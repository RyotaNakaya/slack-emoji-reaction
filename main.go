package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/RyotaNakaya/slack-emoji-reaction/lib"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

var (
	logger    *zap.SugaredLogger
	startTime = *flag.Int("startTime", int(time.Date(2021, 01, 01, 00, 00, 00, 0, time.UTC).Unix()),
		"start unixtime of aggregate")
	endTime = *flag.Int("endTime", int(time.Date(2021, 02, 01, 00, 00, 00, 0, time.UTC).Unix()),
		"end unixtime of aggregate, this is exclusive")
)

func main() {
	logger.Info("start")
	logger.Infof("startTime: %d, endTime: %d", startTime, endTime)
	// lib.GetConversations()
	// lib.GetEmoji()

	aggregateReaction("CF724P8RE", endTime, startTime)

	logger.Info("success!")
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file, error: %v", err)
	}
	lib.SLACK_BOT_TOKEN = os.Getenv("SLACK_BOT_TOKEN")
	lib.SLACK_USER_TOKEN = os.Getenv("SLACK_USER_TOKEN")

	l, err := zap.NewDevelopment()
	if err != nil {
		panic(fmt.Sprintf("failed to create logger: %v", err))
	}
	logger = l.Sugar()

	flag.Parse()
	validateFlags()
}

func validateFlags() {
	if startTime > endTime {
		panic("startTime flag value is late than endTime flag value")
	}
}

func aggregateReaction(ChannelID string, latest int, oldest int) {
	// リアクション
	reactionDict := map[string]int{}
	// スレッドタイムスタンプ
	var ts []string

	// メッセージを取得する
	messages, err := lib.FetchChannelMessages(ChannelID, latest, oldest)
	if err != nil {
		logger.Fatalf("error: %+v", err)
	}
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
	messages, err = lib.FetchChannelThreadMessages("CF724P8RE", ts, latest, oldest)
	if err != nil {
		logger.Fatalf("error: %+v", err)
	}
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
	if err != nil {
		logger.Fatalf("error: %+v", err)
	}
	defer file.Close()

	for k, v := range reactionDict {
		row := fmt.Sprintf("%s: %d\n", k, v)
		_, err := file.WriteString(row)
		if err != nil {
			logger.Fatalf("error: %+v", err)
		}
	}
}
