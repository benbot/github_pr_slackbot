package main

import (
	"encoding/json"
	"github.com/BurntSushi/toml"
	"github.com/nlopes/slack"
	"io/ioutil"
	"net/http"
	"strconv"
)

type Config struct {
	SlackToken string
	Channel    string
}

var (
	config Config
	api    *slack.Client
)

func main() {
	http.HandleFunc("/", http.HandlerFunc(handleGithubHook))
	_, err := toml.DecodeFile("./config.toml", &config)
	if err != nil {
		panic(err)
	}
	api = slack.New(config.SlackToken)

	api.PostMessage(config.Channel, "TEST", slack.PostMessageParameters{})
	http.ListenAndServe(":80", nil)
}

func handleGithubHook(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(500)
	}

	var event PREvent

	err = json.Unmarshal(body, &event)
	if err != nil {
		panic(err)
	}

	if event.Action == "opened" {
		messagePRCreated(event.Number)
	}
}

func messagePRCreated(number int) {
	api.PostMessage(config.Channel, "PR: "+strconv.Itoa(number)+" Created.", slack.PostMessageParameters{})
}
