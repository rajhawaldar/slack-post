package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/acarl005/stripansi"
	"github.com/slack-go/slack"
)

func main() {
	var SLACK_WEBHOOK, message string
	var isPostAll bool

	flag.StringVar(&SLACK_WEBHOOK, "u", "", "Slack webhook URL")
	flag.BoolVar(&isPostAll, "-l", false, "Post message line-by-line")
	flag.Parse()
	if SLACK_WEBHOOK == "" {
		SLACK_WEBHOOK = os.Getenv("SLACK_WEBHOOK_URL")
	}
	if !strings.Contains(SLACK_WEBHOOK, "https://hooks.slack.com") {
		fmt.Fprintf(os.Stderr, "Please set SLACK_WEBHOOK_URL as environment variable or pass it with -u flag.\nRefer: https://api.slack.com/messaging/webhooks#getting_started\n")
		os.Exit(1)
	}
	input, e := os.Stdin.Stat()
	if e != nil {
		log.Panic(e.Error())
	}

	if input.Mode()&os.ModeNamedPipe == 0 {
		os.Exit(0)
	}
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		if isPostAll {
			message += scanner.Text() + "\n"
		} else {
			postSlackMessage(SLACK_WEBHOOK, scanner.Text())
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading: %v", err)
	}
	if isPostAll {
		postSlackMessage(SLACK_WEBHOOK, message)
	}
}

func postSlackMessage(SLACK_WEBHOOK, message string) {
	msg := slack.WebhookMessage{
		Text: stripansi.Strip(message),
	}
	err := slack.PostWebhook(SLACK_WEBHOOK, &msg)
	if err != nil {
		log.Fatal(err)
	}
}
