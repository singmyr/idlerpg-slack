package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

// BaseEvent only contains the type of the event.
// It's used to determine what type of event was sent.
type BaseEvent struct {
	Type    string `json:"type"`
	SubType string `json:"subtype"`
}

// HelloEvent contains the meta data for the "hello" event.
type HelloEvent struct{}

// PongEvent contains the meta data for the p"ong" event.
type PongEvent struct {
	Time    uint64 `json:"time"`
	ReplyTo uint64 `json:"reply_to"`
}

// UserTypingEvent contains the meta data for the "user_typing" event.
type UserTypingEvent struct {
	Channel string `json:"channel"`
	User    string `json:"user"`
}

// MessageEvent contains the meta data for the "message" event.
type MessageEvent struct {
	ClientMessageID      string `json:"client_msg_id"`
	SuppressNotification bool   `json:"suppress_notification"`
	Text                 string `json:"text"`
	User                 string `json:"user"`
	Team                 string `json:"team"`
	UserTeam             string `json:"user_team"`
	SourceTeam           string `json:"source_team"`
	Channel              string `json:"channel"`
	EventTimestamp       string `json:"event_ts"`
	Timestamp            string `json:"ts"`
}

type messageEdited struct {
	User      string `json:"user"`
	Timestamp string `json:"ts"`
}

type message struct {
	ClientMessageID string        `json:"client_msg_id"`
	Text            string        `json:"text"`
	User            string        `json:"user"`
	Team            string        `json:"team"`
	Edited          messageEdited `json:"edited"`
	UserTeam        string        `json:"user_team"`
	SourceTeam      string        `json:"source_team"`
	Channel         string        `json:"channel"`
	Timestamp       string        `json:"ts"`
}

type previousMessage struct {
	ClientMessageID string `json:"client_msg_id"`
	Text            string `json:"text"`
	User            string `json:"user"`
	Team            string `json:"team"`
	Timestamp       string `json:"ts"`
}

// MessageChangedEvent contains the meta data for the event "message_changed" event.
type MessageChangedEvent struct {
	Hidden          bool            `json:"hidden"`
	Message         message         `json:"message"`
	Channel         string          `json:"channel"`
	PreviousMessage previousMessage `json:"previous_message"`
	EventTimestamp  string          `json:"event_ts"`
	Timestamp       string          `json:"ts"`
}

// MessageDeletedEvent contains the meta data for the "message_deleted" event.
type MessageDeletedEvent struct {
	Hidden           bool            `json:"hidden"`
	DeletedTimestamp string          `json:"deleted_ts"`
	Channel          string          `json:"channel"`
	PreviousMessage  previousMessage `json:"previous_message"`
	EventTimestamp   string          `json:"event_ts"`
	Timestamp        string          `json:"ts"`
}

// DesktopNotificationEvent contains the meta data for the "desktop_notification" event.
type DesktopNotificationEvent struct {
	Title           string `json:"title"`
	Subtitle        string `json:"subtitle"`
	Message         string `json:"msg"`
	Timestamp       string `json:"ts"`
	Content         string `json:"content"`
	Channel         string `json:"channel"`
	LaunchURI       string `json:"launchUri"`
	AvatarImage     string `json:"avatarImage"`
	SsbFilename     string `json:"ssbFilename"`
	ImageURI        string `json:"imageUri"`
	IsShared        bool   `json:"is_shared"`
	IsChannelInvite bool   `json:"is_channel_invite"`
	EventTimestamp  string `json:"event_ts"`
}

type item struct {
	Channel   string `json:"channel"`
	Timestamp string `json:"ts"`
}

// ReactionAddedEvent contains the meta data for the "reaction_added" event.
type ReactionAddedEvent struct {
	Item           item   `json:"item"`
	User           string `json:"user"`
	Reaction       string `json:"reaction"`
	ItemUser       string `json:"item_user"`
	EventTimestamp string `json:"event_ts"`
	Timestamp      string `json:"ts"`
}

// ReactionRemovedEvent contains the meta data for the "reaction_removed" event.
type ReactionRemovedEvent struct {
	Item           item   `json:"item"`
	User           string `json:"user"`
	Reaction       string `json:"reaction"`
	ItemUser       string `json:"item_user"`
	EventTimestamp string `json:"event_ts"`
	Timestamp      string `json:"ts"`
}

func connect(token string) {
	const authURL = "https://slack.com/api/rtm.connect?token=%s&presence_sub=1&batch_presence_aware=1"

	response, err := http.Get(fmt.Sprintf(authURL, token))

	if err != nil {
		// handle error
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	type AuthResponse struct {
		Ok  bool   `json:"ok"`
		URL string `json:"url"`
	}
	var v AuthResponse
	err = json.Unmarshal(body, &v)
	if err != nil {
		log.Fatal(err)
	}
	/*{
		"ok": true,
		"self": {
			"id": "      UKWTYQEJV",
			"name": "idle_rpg"
		},
		"team": {
			"domain": "singmyr",
			"id": "T1CAKHC0N",
			"name": "singmyr.io"
		},
		"url": "wss://cerberus-xxxx.lb.slack-msgs.com/websocket/FI7ykBRQQ6b7ylhofO-HcuF0iP6ZmdrNCbM3Dlbdxnkco3TSWmGvYy5ZoF0Upibbps7zmynedGuVnkeuTE_wK3qKw9SwW20aG4q65LBEW0Y="
	}*/
	//data := v.(map[string]interface{})
	/*for k, v := range data {
		switch v := v.(type) {
		case string:
			fmt.Println(k, v, "(string)")
		case float64:
			fmt.Println(k, v, "(float64)")
		case []interface{}:
			fmt.Println(k, "(array):")
			for i, u := range v {
				fmt.Println("    ", i, u)
			}
		default:
			fmt.Println(k, v, "(unknown)")
		}
	}*/

	//wssURL := v.URL
	//fmt.Println("Connecting to ", wssURL)
	//var addr = flag.String("addr", string(wssURL), "http service address")
	flag.Parse()
	log.SetFlags(0)
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	u, err := url.Parse(v.URL)
	if err != nil {
		log.Fatal(err)
	}
	//u := url.URL{Scheme: "ws", Host: *addr, Path: "/echo"}

	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)
			var e BaseEvent
			err = json.Unmarshal(message, &e)
			if err != nil {
				log.Fatal(err)
			}

			err = nil
			switch e.Type {
			case "hello":
				var ev HelloEvent
				err = json.Unmarshal(message, &ev)
				log.Println("Hello")
			case "pong":
				var ev PongEvent
				err = json.Unmarshal(message, &ev)
				log.Println("Pong: ", ev.ReplyTo)
			case "user_typing":
				var ev UserTypingEvent
				err = json.Unmarshal(message, &ev)
				log.Println("User typing:", ev.User, ev.Channel)
			case "message":
				switch e.SubType {
				case "message_changed":
					var ev MessageChangedEvent
					err = json.Unmarshal(message, &ev)
					log.Println("Message changed", ev.Message.Text)
				case "message_deleted":
					var ev MessageDeletedEvent
					err = json.Unmarshal(message, &ev)
					log.Println("Message deleted")
				default:
					var ev MessageEvent
					err = json.Unmarshal(message, &ev)
					log.Println("Message", ev.Text)
				}
			case "desktop_notification":
				var ev DesktopNotificationEvent
				err = json.Unmarshal(message, &ev)
				log.Println("Desktop notification")
			case "reaction_added":
				var ev ReactionAddedEvent
				err = json.Unmarshal(message, &ev)
				log.Println("Reaction added")
			case "reaction_removed":
				var ev ReactionRemovedEvent
				err = json.Unmarshal(message, &ev)
				log.Println("Reaction removed")
			default:
				log.Printf("Unknown event: %v -> %v", e.Type, message)
			}
			if err != nil {
				log.Fatal(err)
			}
		}
	}()

	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			//err := c.WriteMessage(websocket.TextMessage, []byte(t.String()))
			err := c.WriteMessage(websocket.TextMessage, []byte(`{
    "id": 1234,
    "type": "ping",
    "time": 1403299273342
}`))
			if err != nil {
				log.Println("write:", err)
				return
			}
		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}

func main() {
	token := os.Getenv("SLACK_TOKEN")
	connect(token)
}
