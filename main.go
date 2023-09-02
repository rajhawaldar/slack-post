package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	SLACK_WEBHOOK := os.Getenv("SLACK_WEBHOOK_URL")
	var WebHookURL string
	var isPostFile bool
	flag.StringVar(&WebHookURL, "u", "", "Slack webhook URL")
	flag.BoolVar(&isPostFile, "f", false, "Input is a file path")
	flag.Parse()
	tail := flag.Args()

	if !strings.Contains(SLACK_WEBHOOK, "https://hooks.slack.com") {
		if WebHookURL == "" {
			fmt.Fprintf(os.Stderr, "Please set SLACK_WEBHOOK_URL as environment variable or pass it with -u flag")
			os.Exit(1)
		}
	}
	if isPostFile {
		if len(tail) == 0 {
			fmt.Fprintf(os.Stderr, "Please provide file names with -f flag.\n")
			os.Exit(1)
		}
		for _, filePath := range tail {
			if _, err := os.Stat(filePath); err == nil {
				fmt.Println("File Exist:", filePath)
			} else {
				fmt.Fprintf(os.Stderr, filePath+" does not exist\n")

			}
		}
	}
	input, e := os.Stdin.Stat()
	if e != nil {
		log.Panic(e.Error())
	}

	if input.Mode()&os.ModeNamedPipe == 0 {
		os.Exit(0)
	}
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Posting following message on Slack:")
	for scanner.Scan() {
		data := scanner.Text()
		fmt.Println(data)
	}
}
