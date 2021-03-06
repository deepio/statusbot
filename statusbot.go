package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

const (
	UP = iota
	NORMAL
	ERROR
	INFO
	DOWN
)

var colours = [...]string{
	"#36a64f", // Green - Up
	"#ffffff", // White - Normal
	"#db9f49", // Yellow - Error
	"#c0c0c0", // Gray - Info
	"#a50008", // Red - Down
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

type SlackMessage struct {
	Username    string            `json:"username"`
	IconEmoji   string            `json:"icon_emoji"`
	Channel     string            `json:"channel"`
	Text        string            `json:"text,omitempty"`
	Attachments []SlackAttachment `json:"attachments,omitempty"`
}

type SlackAttachment struct {
	Markdown  []string `json:"mrkdwn_in"`
	Colour    string   `json:"color"`
	Title     string   `json:"title,omitempty"`
	Text      string   `json:"text"`
	Footer    string   `json:"footer"`
	TimeStamp int      `json:"ts"`
}

type File struct {
	Sites []Site `yaml:"sites"`
}

type Site struct {
	Name      string `yaml:"name"`
	URL       string `yaml:"url"`
	OldStatus int    `yaml:"omitempty"`
}

func (s Site) GetStatus() (int, error) {
	// r, err := http.Get(s.URL)
	client := http.Client{Timeout: 2 * time.Second}
	r, err := client.Get(s.URL)

	if err != nil {
		return DOWN, err
	}

	if r.StatusCode < 400 {
		return UP, nil
	} else if r.StatusCode >= 400 && r.StatusCode < 500 {
		return DOWN, fmt.Errorf("Status: %d", r.StatusCode)
	} else {
		return ERROR, nil
	}
}

func ParseConf(filepath string) File {
	// If OldStatus is not defined, default is assume its up
	// when statusbot starts because the nil == 0 and UP == 0

	file, _ := ioutil.ReadFile(filepath)
	data := File{}
	_ = yaml.Unmarshal([]byte(file), &data)

	// // Debug
	// // fmt.Printf("%+v", data)
	// for i := 0; i < len(data.Sites); i++ {
	// 	fmt.Printf("Site Name: %-25s\tURL: %s\n", data.Sites[i].Name, data.Sites[i].URL)
	// }

	return data
}

func SlackSend(message string, colour int, channel string) {

	url := getEnv("SLACK_WEBHOOK", "")
	msg := SlackMessage{}
	msg.Username = "StatusBot"
	msg.IconEmoji = ":space_invader:"
	msg.Channel = channel
	msg.Attachments = []SlackAttachment{
		{
			Markdown:  []string{"text"},
			Colour:    colours[colour],
			Text:      message,
			Footer:    "StatusBot - The Newer-er Hotness",
			TimeStamp: int(time.Now().Unix()),
		},
	}

	requestByte, _ := json.Marshal(msg)
	r, err := http.Post(url, "application/json", bytes.NewReader(requestByte))
	if err != nil {
		log.Fatal(err)
	}
	defer r.Body.Close()

	// Debug
	// fmt.Printf("%#v\n", r)
	// body, _ := ioutil.ReadAll(r.Body)
	// fmt.Println(string(body))
}

func Watch(sites File, watch_interval int, channel string) {
	fmt.Println("Starting watch...")

	for {

		for i := 0; i < len(sites.Sites); i++ {
			response, _ := sites.Sites[i].GetStatus()

			if sites.Sites[i].OldStatus == response {
				continue
			} else {
				switch response {
				case UP:
					sites.Sites[i].OldStatus = UP
					fmt.Printf("Site %s is ok.\n", sites.Sites[i].Name)
					msg := fmt.Sprintf("*%s* is back to normal.", sites.Sites[i].Name)
					SlackSend(msg, UP, channel)
				case DOWN:
					fmt.Printf("Site %s is down.\n", sites.Sites[i].Name)
					sites.Sites[i].OldStatus = DOWN
					msg := fmt.Sprintf("*%s* is down!! Link <%s>.", sites.Sites[i].Name, sites.Sites[i].URL)
					SlackSend(msg, DOWN, channel)
				case ERROR:
					fmt.Printf("Site %s is err.\n", sites.Sites[i].Name)
					sites.Sites[i].OldStatus = ERROR
					msg := fmt.Sprintf("*%s* is experiencing errors. Link <%s>.", sites.Sites[i].Name, sites.Sites[i].URL)
					SlackSend(msg, ERROR, channel)
				}
			}
		}

		time.Sleep(time.Second * time.Duration(watch_interval))
	}
}

func main() {
	filepath := flag.String("file", "/tmp/test.yaml", "A filepath to the config file written in YAML or JSON/")
	wait_interval := flag.Int("wait", 15, "Ping interval to use for all sites.")
	channel := flag.String("chan", "#status-bot", "Slack channel to send notifications to.")
	flag.Parse()

	// Check if Environment Variable is set
	url := getEnv("SLACK_WEBHOOK", "")
	if url == "" {
		fmt.Println(`
The environment variable "SLACK_WEBHOOK" does not exist.
Please add the slack webhook environment variable to use this software.

export SLACK_WEBHOOK=https://hooks.slack.com/services/XXXXXXXXX/XXXXXXXXX/XXXXXXXXXXXXXXXXXXXXXXXX
		`)
		os.Exit(1)
	}

	SlackSend("Statusbot is connected.", NORMAL, *channel)
	f := ParseConf(*filepath)
	Watch(f, *wait_interval, *channel)
}
