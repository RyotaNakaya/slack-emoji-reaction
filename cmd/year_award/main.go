package main

import (
	"fmt"
	"os"
	"time"

	"github.com/namsral/flag"

	"github.com/RyotaNakaya/slack-emoji-reaction/lib"
	"github.com/RyotaNakaya/slack-emoji-reaction/lib/repository"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

var (
	now         = time.Now()
	logger      *zap.SugaredLogger
	slackClient *lib.Slack

	// test 用
	targetChannelID = flag.String("targetChannelID", "C023TM73HR7", "target post shannel_id")
	startTime       = flag.Int("startTime", int(time.Date(now.Year(), 1, 1, 0, 0, 0, 0, time.Local).Unix()),
		"start unixtime of aggregate")
	endTime = flag.Int("endTime", int(time.Date(now.Year(), 12, 31, 23, 59, 59, 99, time.Local).Unix()),
		"end unixtime of aggregate, this is exclusive")

	dbUser         = flag.String("dbuser", "root", "mysql user name")
	dbPass         = flag.String("dbpass", "", "mysql password")
	dbHost         = flag.String("dbhost", "localhost", "mysql host name")
	dbPort         = flag.String("dbport", "3306", "mysql port")
	dbName         = flag.String("dbname", "slack_reaction_development", "mysql database name")
	dbMaxOpenConn  = flag.Int("dbcon_open", 100, "mysql connection pool's max open connection quantity")
	dbMaxIdleConn  = flag.Int("dbcon_idle", 100, "mysql connection pool's max idle connection quantity")
	dbConnLifetime = flag.String("dbcon_timeout", "1h", "mysql connection's lifetime in seconds")
)

func main() {
	defer func() { _ = logger.Sync() }()
	defer repository.DB.Close()

	logger.Info("start")

	slackClient = lib.NewSlack(os.Getenv("SLACK_BOT_TOKEN"))

	// リアクション数が多かったアワード
	records := aggregateReactionCountAward(*startTime, *endTime, 3)
	header := `
#####################################################
:tada: :tada: *リアクション数が多かったアワード* :tada: :tada:
#####################################################
	`
	if err := slackClient.PostMessage(*targetChannelID, header); err != nil {
		logger.Fatalf("error: %+v", err)
	}
	for k, record := range records {
		text := fmt.Sprintf("%d位 (%d 個)\n %s", k+1, record.count, record.link)
		if err := slackClient.PostMessage(*targetChannelID, text); err != nil {
			logger.Fatalf("error: %+v", err)
		}
	}

	// リアクション種類数が多かったアワード
	kindRecords := aggregateReactionKindCountAward(*startTime, *endTime, 3)
	header = `
#####################################################
:tada: :tada: *リアクション種類数が多かったアワード* :tada: :tada:
#####################################################
		`
	if err := slackClient.PostMessage(*targetChannelID, header); err != nil {
		logger.Fatalf("error: %+v", err)
	}
	for k, record := range kindRecords {
		text := fmt.Sprintf("%d位 (%d 個)\n %s", k+1, record.count, record.link)
		if err := slackClient.PostMessage(*targetChannelID, text); err != nil {
			logger.Fatalf("error: %+v", err)
		}
	}

	// omoro アワード
	omoroRecords := aggregateReactionOmoroCountAward(*startTime, *endTime, 3)
	header = `
#####################################################
:tada: :tada: *オモロアワード* :tada: :tada: :wwww: :kusa: :omoroi: :warota: :kusa_1:
#####################################################
	`
	if err := slackClient.PostMessage(*targetChannelID, header); err != nil {
		logger.Fatalf("error: %+v", err)
	}
	for k, record := range omoroRecords {
		text := fmt.Sprintf("%d位 (%d 個)\n %s", k+1, record.count, record.link)
		if err := slackClient.PostMessage(*targetChannelID, text); err != nil {
			logger.Fatalf("error: %+v", err)
		}
	}

	// good story アワード
	goodStoryRecords := aggregateReactionGoodStoryCountAward(*startTime, *endTime, 3)
	header = `
#####################################################
:tada: :tada: *いい話アワード* :tada: :tada: :iihanasi: :iihanasi: :iihanasi:
#####################################################
	`
	if err := slackClient.PostMessage(*targetChannelID, header); err != nil {
		logger.Fatalf("error: %+v", err)
	}
	for k, record := range goodStoryRecords {
		text := fmt.Sprintf("%d位 (%d 個)\n %s", k+1, record.count, record.link)
		if err := slackClient.PostMessage(*targetChannelID, text); err != nil {
			logger.Fatalf("error: %+v", err)
		}
	}

	logger.Info("success!")
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

	// parse flag
	flag.Parse()

	// setup database
	repository.DB = repository.PrepareDBConnection(*dbUser, *dbPass, *dbHost, *dbPort, *dbName, *dbMaxOpenConn, *dbMaxIdleConn, *dbConnLifetime)
}

type record struct {
	link  string
	count int
}

// リアクション数が多いメッセージのリンクを limit 分だけ返します
func aggregateReactionCountAward(start, end, limit int) []record {
	q := `select channel_id, message_ts_nano, sum(reaction_count) count
		from message_reactions
		where message_ts between ? and ?
		group by channel_id, message_ts_nano
		order by sum(reaction_count) desc
		limit ?;`

	var res []reactionCountAward
	if err := repository.DB.Select(&res, q, start, end, limit); err != nil {
		logger.Fatalf("error: %+v", err)
	}

	records := []record{}
	for _, v := range res {
		link, err := slackClient.FetchMessage(v.ChannelID, v.MessageTsNano)
		if err != nil {
			logger.Fatalf("error: %+v", err)
		}
		record := record{
			link:  link,
			count: v.Count,
		}
		records = append(records, record)
	}

	return records
}

type reactionCountAward struct {
	ChannelID     string `db:"channel_id"`
	MessageTsNano string `db:"message_ts_nano"`
	Count         int    `db:"count"`
}

// リアクションの種類数が多いメッセージのリンクを limit 分だけ返します
func aggregateReactionKindCountAward(start, end, limit int) []record {
	q := `select channel_id, message_ts_nano, count(message_ts_nano) count
		from message_reactions
		where message_ts between ? and ?
		group by channel_id, message_ts_nano
		order by count(message_ts_nano) desc
		limit ?;`

	var res []reactionKindCountAward
	if err := repository.DB.Select(&res, q, start, end, limit); err != nil {
		logger.Fatalf("error: %+v", err)
	}

	records := []record{}
	for _, v := range res {
		link, err := slackClient.FetchMessage(v.ChannelID, v.MessageTsNano)
		if err != nil {
			logger.Fatalf("error: %+v", err)
		}
		record := record{
			link:  link,
			count: v.Count,
		}
		records = append(records, record)
	}

	return records
}

type reactionKindCountAward struct {
	ChannelID     string `db:"channel_id"`
	MessageTsNano string `db:"message_ts_nano"`
	Count         int    `db:"count"`
}

// omoro 系リアクションの数が多いメッセージのリンクを limit 分だけ返します
func aggregateReactionOmoroCountAward(start, end, limit int) []record {
	q := `select channel_id, message_ts_nano, sum(reaction_count) count
		from message_reactions
		where message_ts between ? and ?
		and reaction_name in("wwww", "kusa", "kusa_1", "omoroi", "warota")
		group by channel_id, message_ts_nano
		order by sum(reaction_count) desc
		limit ?;`

	var res []reactionOmoroCountAward
	if err := repository.DB.Select(&res, q, start, end, limit); err != nil {
		logger.Fatalf("error: %+v", err)
	}

	records := []record{}
	for _, v := range res {
		link, err := slackClient.FetchMessage(v.ChannelID, v.MessageTsNano)
		if err != nil {
			logger.Fatalf("error: %+v", err)
		}
		record := record{
			link:  link,
			count: v.Count,
		}
		records = append(records, record)
	}

	return records
}

type reactionOmoroCountAward struct {
	ChannelID     string `db:"channel_id"`
	MessageTsNano string `db:"message_ts_nano"`
	Count         int    `db:"count"`
}

// iihanashi リアクションの数が多いメッセージのリンクを limit 分だけ返します
func aggregateReactionGoodStoryCountAward(start, end, limit int) []record {
	q := `select channel_id, message_ts_nano, sum(reaction_count) count
		from message_reactions
		where message_ts between ? and ?
		and reaction_name in("iihanasi")
		group by channel_id, message_ts_nano
		order by sum(reaction_count) desc
		limit ?;`

	var res []reactionGoodStoryCountAward
	if err := repository.DB.Select(&res, q, start, end, limit); err != nil {
		logger.Fatalf("error: %+v", err)
	}

	records := []record{}
	for _, v := range res {
		link, err := slackClient.FetchMessage(v.ChannelID, v.MessageTsNano)
		if err != nil {
			logger.Fatalf("error: %+v", err)
		}
		record := record{
			link:  link,
			count: v.Count,
		}
		records = append(records, record)
	}

	return records
}

type reactionGoodStoryCountAward struct {
	ChannelID     string `db:"channel_id"`
	MessageTsNano string `db:"message_ts_nano"`
	Count         int    `db:"count"`
}
