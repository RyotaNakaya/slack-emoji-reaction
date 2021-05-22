package lib

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/slack-go/slack"
)

func GetChannelHistory(ids []string, latest int, oldest int) {
	api := slack.New(SLACK_USER_TOKEN)
	param := slack.GetConversationHistoryParameters{
		ChannelID: ids[0],
		Cursor:    "",
		Inclusive: false,
		Latest:    strconv.Itoa(int(latest)),
		Limit:     200,
		Oldest:    strconv.Itoa(int(oldest)),
	}

	// next cursor が返ってこなくなるまで再帰的にコール
	for {
		res, err := api.GetConversationHistory(&param)
		if err != nil {
			log.Fatal(err)
		}

		output := make([]map[string]string, 0, len(res.Messages))
		for _, message := range res.Messages {
			m := make(map[string]string)
			m["text"] = message.Msg.Text
			m["time"] = message.Msg.Timestamp
			output = append(output, m)
		}

		var filename string = "tmp/msg.txt"
		file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
		defer file.Close()
		if err != nil {
			log.Fatal(err)
		}

		for _, v := range output {
			s := ""
			for k, v := range v {
				s += fmt.Sprintf("%s: %s,", k, v)
			}
			s += "\n"
			_, err := file.WriteString(s)
			if err != nil {
				log.Fatal(err)
			}
		}

		// 次のページがなければ終了
		next := res.ResponseMetaData.NextCursor
		if next == "" {
			break
		} else {
			param.Cursor = next
		}
	}
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
