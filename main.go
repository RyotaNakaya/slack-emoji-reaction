package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/RyotaNakaya/slack-emoji-reaction/lib"
	"github.com/RyotaNakaya/slack-emoji-reaction/lib/repository"
	"github.com/joho/godotenv"
	"github.com/slack-go/slack"
	"go.uber.org/zap"
)

var (
	now    = time.Now()
	logger *zap.SugaredLogger

	targetChannelID = *flag.String("targetChannelID", "", "if not specify, aggregate all public ch")
	startTime       = *flag.Int("startTime", int(time.Date(now.Year(), now.Month()-1, 1, 0, 0, 0, 0, time.Local).Unix()),
		"start unixtime of aggregate")
	endTime = *flag.Int("endTime", int(time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local).Unix()),
		"end unixtime of aggregate, this is exclusive")

	dbUser         = *flag.String("dbuser", "root", "mysql user name")
	dbPass         = *flag.String("dbpass", "", "mysql password")
	dbHost         = *flag.String("dbhost", "localhost", "mysql host name")
	dbPort         = *flag.String("dbport", "3306", "mysql port")
	dbName         = *flag.String("dbname", "slack_reaction_development", "mysql database name")
	dbMaxOpenConn  = *flag.Int("dbcon_open", 100, "mysql connection pool's max open connection quantity")
	dbMaxIdleConn  = *flag.Int("dbcon_idle", 100, "mysql connection pool's max idle connection quantity")
	dbConnLifetime = *flag.String("dbcon_timeout", "1h", "mysql connection's lifetime in seconds")
)

func main() {
	defer func() { _ = logger.Sync() }()

	st := time.Now()
	logger.Info("start")
	logger.Infof("startTime: %d(%v), endTime: %d(%v)", startTime, time.Unix(int64(startTime), 0), endTime, time.Unix(int64(endTime), 0))

	// 集計対象のチャンネルを取得
	var chs []string
	var err error
	if targetChannelID == "" {
		chs, err = lib.FetchPublicChannelIDs()
		if err != nil {
			logger.Fatalf("error: %+v", err)
		}
	} else {
		chs = []string{targetChannelID}

	}

	// 時間かかるから...一旦絞る
	if len(chs) > 5 {
		chs = chs[0:4]
	}

	// reaction の集計
	for idx, cid := range chs {
		logger.Infof("aggregate ch: %s, idx: %d", cid, idx)
		aggregateReaction(cid, endTime, startTime)
	}

	logger.Info("success!")
	et := time.Now()
	logger.Infof("The call took %v to run.\n", et.Sub(st))

}

func init() {
	// setup logger
	l, err := zap.NewDevelopment()
	if err != nil {
		panic(fmt.Sprintf("failed to create logger: %v", err))
	}
	logger = l.Sugar()

	// read env variable
	if err := godotenv.Load(); err != nil {
		logger.Fatalf("Error loading .env file, error: %v", err)
	}
	lib.SLACK_BOT_TOKEN = os.Getenv("SLACK_BOT_TOKEN")
	lib.SLACK_USER_TOKEN = os.Getenv("SLACK_USER_TOKEN")

	// parse flag
	flag.Parse()
	validateFlags()

	// setup database
	repository.DB = repository.PrepareDBConnection(dbUser, dbPass, dbHost, dbPort, dbName, dbMaxOpenConn, dbMaxIdleConn, dbConnLifetime)
}

func validateFlags() {
	if startTime > endTime {
		panic("startTime flag value is late than endTime flag value")
	}
}

func aggregateReaction(ChannelID string, latest int, oldest int) {
	logger.Infof("aggregate channel: %s", ChannelID)
	mrs := repository.MessageReactions{}
	now := time.Now().Unix()
	// スレッドタイムスタンプ
	var ts []string

	// メッセージを取得する
	messages, err := lib.FetchChannelMessages(ChannelID, latest, oldest)
	if err != nil {
		logger.Fatalf("error: %+v", err)
	}
	// reaction 集計
	mrs.MessageReactions = append(mrs.MessageReactions, buildMessageReactions(ChannelID, now, messages)...)
	// TODO: 無駄に二周ループ回してるのでもっとうまくできそう
	for _, message := range messages {
		// ThreadTimestamp がある場合スレッド取得で使うので溜めておく
		if t := message.Msg.ThreadTimestamp; t != "" {
			ts = append(ts, t)
		}
	}

	// スレッドを取得する
	logger.Infof("start get thread: %#v", ts)
	messages, err = lib.FetchChannelThreadMessages(ChannelID, ts, latest, oldest)
	if err != nil {
		logger.Fatalf("error: %+v", err)
	}
	mrs.MessageReactions = append(mrs.MessageReactions, buildMessageReactions(ChannelID, now, messages)...)

	if err = mrs.Save(); err != nil {
		logger.Fatalf("error: %+v, messages: %#v", err, messages)
	}
}

func buildMessageReactions(chid string, now int64, messages []slack.Message) []*repository.MessageReaction {
	res := []*repository.MessageReaction{}
	for _, message := range messages {
		tsUnix, err := strconv.Atoi(message.Timestamp[0:10])
		if err != nil {
			logger.Fatalf("error: %+v", err)
		}

		t := time.Unix(int64(tsUnix), 0)

		for _, r := range message.Reactions {
			mr := repository.MessageReaction{
				ChannelID:     chid,
				MessageID:     message.Msg.ClientMsgID,
				ReactionName:  r.Name,
				ReactionCount: uint(r.Count),
				MessageTS:     message.Timestamp,
				YYYYMM:        strconv.Itoa(t.Year()) + fmt.Sprintf("%02d", int(t.Month())),
				CreatedAt:     uint(now),
			}
			res = append(res, &mr)
		}
	}
	return res
}
