package lib

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func GetConversations() {
	var next string
	client := new(http.Client)
	// 再帰的にページングしてリクエストを送る
	for {
		req, _ := http.NewRequest("GET", buildConversationsUrl(next), nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", SLACK_TOKEN))

		res, _ := client.Do(req)
		body, _ := ioutil.ReadAll(res.Body)
		var cr ChannelResponse
		if err := json.Unmarshal(body, &cr); err != nil {
			log.Fatal(err)
		}
		outputConversationsData(cr.Channels)
		next = cr.Metadata["next_cursor"]
		// 次のページがなければ終了
		if next == "" {
			break
		}
	}
}

func buildConversationsUrl(cursor string) string {
	return fmt.Sprintf("%s//conversations.list?exclude_members=true&exclude_archived=true&limit=200&cursor=%s", BASE_URL, cursor)
}

func outputConversationsData(c []ChannelResponseDetail) {
	outputToFile(c)
}

func outputToFile(c []ChannelResponseDetail) {
	var filename string = "tmp/channels.txt"
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	for _, v := range c {
		_, err := file.WriteString(fmt.Sprintf("%s\n", v.Name))
		if err != nil {
			log.Fatal(err)
		}
	}
}

type ChannelResponse struct {
	Channels []ChannelResponseDetail `json:"channels"`
	Metadata map[string]string       `json:"response_metadata"`
}

type ChannelResponseDetail struct {
	Name string `json:"name"`
	Id   string `json:"id"`
}
