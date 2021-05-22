package lib

import (
	"fmt"
	"log"
	"os"

	"github.com/slack-go/slack"
)

func GetEmoji() {
	api := slack.New(SLACK_BOT_TOKEN)
	res, err := api.GetEmoji()
	if err != nil {
		log.Fatal(err)
	}

	var filename string = "tmp/emoji.txt"
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	for k := range res {
		_, err := file.WriteString(fmt.Sprintf("%s\n", k))
		if err != nil {
			log.Fatal(err)
		}
	}
}
