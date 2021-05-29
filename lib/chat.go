package lib

import (
	"github.com/slack-go/slack"
)

func (s *Slack) PostMessage(chid string, text string) error {
	t := slack.MsgOptionText(text, false)
	_, _, err := s.client.PostMessage(chid, t)
	if err != nil {
		return err
	}
	return nil
}
