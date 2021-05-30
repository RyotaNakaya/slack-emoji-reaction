package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/RyotaNakaya/slack-emoji-reaction/lib"
	"github.com/RyotaNakaya/slack-emoji-reaction/lib/repository"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

var (
	now    = time.Now()
	logger *zap.SugaredLogger

	targetChannelID = *flag.String("targetChannelID", "C023TM73HR7", "target post shannel_id")
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
	logger.Info("start post_reaction_result")
	logger.Infof("startTime: %d(%v), endTime: %d(%v)", startTime, time.Unix(int64(startTime), 0), endTime, time.Unix(int64(endTime), 0))

	ms := selectReactionCount(startTime, endTime)
	text := "今月の slack リアクション集計結果"
	for i, v := range ms {
		text += fmt.Sprintf("\n%d位: :%s: %d回", i+1, v.ReactionName, v.ReactionCount)
	}

	s := lib.NewSlack(os.Getenv("SLACK_BOT_TOKEN"))
	err := s.PostMessage(targetChannelID, text)
	if err != nil {
		logger.Fatalf("error: %+v", err)
	}

	logger.Info("success post_reaction_result!")
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

// selectReactionCount はリアクション数の多い順にソートした結果を返します
func selectReactionCount(s, e int) []messageReaction {
	res := []messageReaction{}
	q := `
		select reaction_name, sum(reaction_count) reaction_count
		from message_reactions
		where message_ts between ? and ?
		group by reaction_name
		order by reaction_count desc
		limit 10;`

	rows, err := repository.DB.Queryx(q, s, e)
	if err != nil {
		logger.Fatalf("error: %+v", err)
	}
	var r messageReaction
	for rows.Next() {
		err = rows.StructScan(&r)
		res = append(res, r)
	}
	if err != nil {
		logger.Fatalf("error: %+v", err)
	}
	return res
}

type messageReaction struct {
	ReactionName  string `db:"reaction_name"`
	ReactionCount int    `db:"reaction_count"`
}
