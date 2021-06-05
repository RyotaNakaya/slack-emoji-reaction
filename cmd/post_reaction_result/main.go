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
	now    = time.Now()
	logger *zap.SugaredLogger

	// test 用
	targetChannelID = flag.String("targetChannelID", "C023TM73HR7", "target post shannel_id")
	// 本番用
	// targetChannelID = flag.String("targetChannelID", "C023SG46EBF", "target post shannel_id")
	startTime = flag.Int("startTime", int(time.Date(now.Year(), now.Month()-1, 1, 0, 0, 0, 0, time.Local).Unix()),
		"start unixtime of aggregate")
	endTime = flag.Int("endTime", int(time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local).Unix()),
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

const (
	MOST_USED_REACTION_POST_COUNT     = 10
	MOST_USED_REACTED_USER_POST_COUNT = 10
)

func main() {
	defer func() { _ = logger.Sync() }()
	defer repository.DB.Close()

	st := time.Now()
	logger.Info("start post_reaction_result")
	logger.Infof("startTime: %d(%v), endTime: %d(%v)", *startTime, time.Unix(int64(*startTime), 0), *endTime, time.Unix(int64(*endTime), 0))

	term := fmt.Sprintf("[%v~%v]", time.Unix(int64(*startTime), 0).Format("2006/01/02"), time.Unix(int64(*endTime), 0).Format("2006/01/02"))
	// 多く使われたリアクションを集計
	ms := selectReactionCount(*startTime, *endTime)
	text := ":tada: *たくさん使われたリアクション* :tada: " + term
	for i, v := range ms {
		text += fmt.Sprintf("\n%d位: :%s: %d回", i+1, v.ReactionName, v.ReactionCount)
	}

	// リアクションをたくさんもらった人を集計
	// TODO: #general や #kintai のような業務連絡系は除外してもいいかも
	text += "\n\n:trophy: *リアクションをたくさんもらった人* :trophy: " + term
	ru := selectReactedUser(*startTime, *endTime)
	for _, v := range ru {
		t := ""
		for _, v2 := range v.ReactedUserReaction {
			t += fmt.Sprintf(":%s: %d回、", v2.ReactionName, v2.ReactionCount)
		}
		text += fmt.Sprintf("\n%s さん=> %s...etc", v.UserName, t)
	}

	// slack にポストする
	s := lib.NewSlack(os.Getenv("SLACK_BOT_TOKEN"))
	err := s.PostMessage(*targetChannelID, text)
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
	repository.DB = repository.PrepareDBConnection(*dbUser, *dbPass, *dbHost, *dbPort, *dbName, *dbMaxOpenConn, *dbMaxIdleConn, *dbConnLifetime)
}

func validateFlags() {
	if *startTime > *endTime {
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
		limit ?;`

	rows, err := repository.DB.Queryx(q, s, e, MOST_USED_REACTION_POST_COUNT)
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

// selectReactedUser はリアクションをたくさんもらったユーザーを返します
// TODO: ハイパーやっつけ
func selectReactedUser(s, e int) []ReactedUser {
	res := []ReactedUser{}
	// TODO: Slackbot を弾く
	q := `
		select message_user_id, sum(reaction_count) reaction_count
		from message_reactions
		where message_ts between ? and ?
		and message_user_id != ""
		group by message_user_id
		order by reaction_count desc
		limit ?;`

	rows, err := repository.DB.Queryx(q, s, e, MOST_USED_REACTED_USER_POST_COUNT)
	if err != nil {
		logger.Fatalf("error: %+v", err)
	}
	userIDs := []string{}
	for rows.Next() {
		var message_user_id string
		var reaction_count int
		err = rows.Scan(&message_user_id, &reaction_count)
		userIDs = append(userIDs, message_user_id)
	}
	if err != nil {
		logger.Fatalf("error: %+v", err)
	}

	c := lib.NewSlack(os.Getenv("SLACK_USER_TOKEN"))
	for _, v := range userIDs {
		ru := ReactedUser{}
		u, err := c.GetUser(v)
		if err != nil {
			logger.Fatalf("error: %+v", err)
		}
		ru.UserName = u.DisplayName

		q = `select reaction_name, sum(reaction_count) reaction_count
			from message_reactions
			where message_ts between ? and ?
			and message_user_id = ?
			group by reaction_name
			order by reaction_count desc
			limit 5;`
		rows, err := repository.DB.Queryx(q, s, e, v)
		if err != nil {
			logger.Fatalf("error: %+v", err)
		}
		for rows.Next() {
			rur := ReactedUserReaction{}
			var reaction_name string
			var reaction_count int
			err = rows.Scan(&reaction_name, &reaction_count)
			rur.ReactionName = reaction_name
			rur.ReactionCount = reaction_count
			ru.ReactedUserReaction = append(ru.ReactedUserReaction, rur)
		}
		if err != nil {
			logger.Fatalf("error: %+v", err)
		}
		res = append(res, ru)
	}
	return res
}

type ReactedUser struct {
	UserName            string
	ReactedUserReaction []ReactedUserReaction
}

type ReactedUserReaction struct {
	ReactionName  string
	ReactionCount int
}
