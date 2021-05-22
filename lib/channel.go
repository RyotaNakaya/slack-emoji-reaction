package lib

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/slack-go/slack"
)

func GetChannelHistory(ChannelID string, latest int, oldest int) []slack.Message {
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
		r, err := api.GetConversationHistory(&param)
		if err != nil {
			log.Fatal(err)
		}

		for _, v := range r.Messages {
			res = append(res, v)
		}

		// 次のページがなければ終了
		next := r.ResponseMetaData.NextCursor
		if next == "" {
			break
		} else {
			param.Cursor = next
		}
	}

	return res
}

func GetConversations() {
	api := slack.New(SLACK_BOT_TOKEN)
	param := slack.GetConversationsParameters{
		ExcludeArchived: false,
		Limit:           200,
		Types:           []string{"public_channel"},
	}

	// next cursor が返ってこなくなるまで再帰的にコール
	for {
		res, next, err := api.GetConversations(&param)
		if err != nil {
			log.Fatal(err)
		}

		// チャンネルIDとチャンネル名を書き出す
		m := make(map[string]string, len(res))
		for _, v := range res {
			m[v.ID] = v.Name
		}
		outputChannelInfoToFile(m)

		// 次のページがなければ終了
		if next == "" {
			break
		} else {
			param.Cursor = next
		}
	}
}

func outputChannelInfoToFile(c map[string]string) {
	var filename string = "tmp/channels.txt"
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	for k, v := range c {
		_, err := file.WriteString(fmt.Sprintf("%s: %s\n", k, v))
		if err != nil {
			log.Fatal(err)
		}
	}
}
