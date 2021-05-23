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
	st := time.Now()
	logger.Info("start")
	logger.Infof("startTime: %d, endTime: %d", startTime, endTime)

	// 集計対象のチャンネルを取得
	chs, err := lib.FetchPublicChannelIDs()
	if err != nil {
		logger.Fatalf("error: %+v", err)
	}

	// 時間かかるから...一旦100個までに絞る
	chs = chs[0:100]

	// reaction の集計
	// リアクション
	reactionDict := map[string]int{}
	for idx, cid := range chs {
		logger.Infof("aggregate ch: %s, idx: %d", cid, idx)
		reactionDict = aggregateReaction(cid, endTime, startTime, reactionDict)
	}

	if err = output(reactionDict); err != nil {
		logger.Fatalf("error: %+v", err)
	}

	logger.Info("success!")
	et := time.Now()
	logger.Infof("The call took %v to run.\n", et.Sub(st))

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

func aggregateReaction(ChannelID string, latest int, oldest int, reactionDict map[string]int) map[string]int {
	// リアクション
	// reactionDict := map[string]int{}
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
	messages, err = lib.FetchChannelThreadMessages(ChannelID, ts, latest, oldest)
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

	return reactionDict
}

func output(reactionDict map[string]int) error {
	var filename string = "tmp/msg.txt"
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("failed to OpenFile: %w", err)
	}
	defer file.Close()

	for k, v := range reactionDict {
		row := fmt.Sprintf("%s: %d\n", k, v)
		_, err := file.WriteString(row)
		if err != nil {
			return fmt.Errorf("failed to WriteString: %w", err)
		}
	}

	return nil
}
