package lib

import (
	"fmt"
	"strconv"
	"time"

	"github.com/slack-go/slack"
)

const retlyCount = 3

// 指定されたチャンネル、期間のメッセージ一覧を返す
func (s *Slack) FetchChannelMessages(ChannelID string, latest int, oldest int) ([]slack.Message, error) {
	param := slack.GetConversationHistoryParameters{
		ChannelID: ChannelID,
		Cursor:    "",
		Inclusive: false,
		Latest:    strconv.Itoa(int(latest)),
		Limit:     500,
		Oldest:    strconv.Itoa(int(oldest)),
	}

	res := []slack.Message{}

	// next cursor が返ってこなくなるまで再帰的にコール
	for {
		// 1分あたり50まで
		// rate limit に引っかからないようにゆっくり叩く
		time.Sleep(time.Second * 1)

		errCount := 0
		var r *slack.GetConversationHistoryResponse
		var err error
		for errCount < retlyCount {
			r, err = s.client.GetConversationHistory(&param)
			if err != nil {
				errCount++
			} else {
				break
			}
		}
		if errCount == retlyCount {
			return nil, fmt.Errorf("failed to GetConversationHistory: %w", err)
		}

		res = append(res, r.Messages...)

		// 次のページがなければ終了
		next := r.ResponseMetaData.NextCursor
		if next == "" {
			break
		} else {
			param.Cursor = next
		}
	}

	return res, nil
}

// 指定されたチャンネルスレッド、期間のメッセージ一覧を返す
func (s *Slack) FetchChannelThreadMessages(ChannelID string, timestamps []string, latest int, oldest int) ([]slack.Message, error) {
	res := []slack.Message{}

	for _, ts := range timestamps {
		// rate limit に引っかからないようにゆっくり叩く
		// 1分あたり50まで
		time.Sleep(time.Second * 1)
		param := slack.GetConversationRepliesParameters{
			ChannelID: ChannelID,
			Timestamp: ts,
			Cursor:    "",
			Inclusive: false,
			Latest:    strconv.Itoa(int(latest)),
			Limit:     500,
			Oldest:    strconv.Itoa(int(oldest)),
		}

		// next cursor が返ってこなくなるまで再帰的にコール
		for {
			errCount := 0
			var r []slack.Message
			var err error
			var next string
			for errCount < retlyCount {
				r, _, next, err = s.client.GetConversationReplies(&param)
				if err != nil {
					errCount++
				} else {
					break
				}
			}
			if errCount == retlyCount {
				return nil, fmt.Errorf("failed to GetConversationReplies: %w", err)
			}

			// 先頭は ConversationHistory で取得済みなので弾く
			res = append(res, r[1:]...)

			// 次のページがなければ終了
			if next == "" {
				break
			} else {
				param.Cursor = next
			}
		}
	}

	return res, nil
}

// パブリックチャンネルIDの一覧を返します
func (s *Slack) FetchPublicChannelIDs() ([]string, error) {
	param := slack.GetConversationsParameters{
		ExcludeArchived: false,
		Limit:           500,
		Types:           []string{"public_channel"},
	}

	res := []string{}

	// next cursor が返ってこなくなるまで再帰的にコール
	for {
		r, next, err := s.client.GetConversations(&param)
		if err != nil {
			return nil, fmt.Errorf("failed to FetchPublicChannelIDs: %w", err)
		}

		for _, v := range r {
			res = append(res, v.ID)
		}

		// 次のページがなければ終了
		if next == "" {
			break
		} else {
			param.Cursor = next
		}
	}

	return res, nil
}
