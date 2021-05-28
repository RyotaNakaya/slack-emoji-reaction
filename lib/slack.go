package lib

import "github.com/slack-go/slack"

var (
	SLACK_BOT_TOKEN  string
	SLACK_USER_TOKEN string
	BASE_URL         = "https://slack.com/api/"
)

type Slack struct {
	client *slack.Client
}

func NewSlack(token string) *Slack {
	s := new(Slack)
	s.client = slack.New(token)
	return s
}
