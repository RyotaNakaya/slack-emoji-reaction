package main

import (
	"fmt"
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
