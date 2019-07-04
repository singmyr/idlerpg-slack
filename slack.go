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
