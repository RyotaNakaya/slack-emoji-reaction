package lib

import (
	"github.com/slack-go/slack"
)

func (s *Slack) FetchMessage(channelID, messageTS string) (string, error) {
	param := slack.PermalinkParameters{
		Channel: channelID,
		Ts:      messageTS,
	}

	return s.client.GetPermalink(&param)
}
