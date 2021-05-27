package lib

import (
	"fmt"
	"strconv"
	"time"

	"github.com/slack-go/slack"
)

// 指定されたチャンネル、期間のメッセージ一覧を返す
func FetchChannelMessages(ChannelID string, latest int, oldest int) ([]slack.Message, error) {
	api := slack.New(SLACK_USER_TOKEN)
	param := slack.GetConversationHistoryParameters{
		ChannelID: ChannelID,
		Cursor:    "",
		Inclusive: false,
		Latest:    strconv.Itoa(int(latest)),
		Limit:     200,
		Oldest:    strconv.Itoa(int(oldest)),
	}

	res := []slack.Message{}

	// next cursor が返ってこなくなるまで再帰的にコール
	for {
		// 1分あたり50まで
		// rate limit に引っかからないようにゆっくり叩く
		time.Sleep(time.Second * 1)
		r, err := api.GetConversationHistory(&param)
		if err != nil {
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
func FetchChannelThreadMessages(ChannelID string, timestamps []string, latest int, oldest int) ([]slack.Message, error) {
	api := slack.New(SLACK_USER_TOKEN)
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
			Limit:     200,
			Oldest:    strconv.Itoa(int(oldest)),
		}

		// next cursor が返ってこなくなるまで再帰的にコール
		for {
			r, _, next, err := api.GetConversationReplies(&param)
			if err != nil {
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
func FetchPublicChannelIDs() ([]string, error) {
	api := slack.New(SLACK_USER_TOKEN)
	param := slack.GetConversationsParameters{
		ExcludeArchived: false,
		Limit:           200,
		Types:           []string{"public_channel"},
	}

	res := []string{}

	// next cursor が返ってこなくなるまで再帰的にコール
	for {
		r, next, err := api.GetConversations(&param)
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
