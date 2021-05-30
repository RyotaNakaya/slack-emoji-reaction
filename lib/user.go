package lib

import "github.com/slack-go/slack"

func (s *Slack) GetUser(uid string) (*slack.UserProfile, error) {
	u, err := s.client.GetUserProfile(&slack.GetUserProfileParameters{UserID: uid})
	if err != nil {
		return nil, err
	}
	return u, nil
}
