package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/acarl005/stripansi"
	"github.com/slack-go/slack"
)

var wg sync.WaitGroup

func main() {
	var SLACK_WEBHOOK, message string
	var postLineByLine bool

	flag.StringVar(&SLACK_WEBHOOK, "u", "", "Slack webhook URL")
	flag.BoolVar(&postLineByLine, "l", false, "Post message line-by-line")
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
		if postLineByLine {
			wg.Add(1)
			go postSlackMessage(SLACK_WEBHOOK, scanner.Text())
		} else {
			message += scanner.Text() + "\n"
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading: %v", err)
	}
	if !postLineByLine {
		wg.Add(1)
		go postSlackMessage(SLACK_WEBHOOK, message)
	}
	wg.Wait()
}

func postSlackMessage(SLACK_WEBHOOK, message string) {
	msg := slack.WebhookMessage{
		Text: stripansi.Strip(message),
	}
	err := slack.PostWebhook(SLACK_WEBHOOK, &msg)
	if err != nil {
		log.Fatal(err)
	}
	defer wg.Done()
}
