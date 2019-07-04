package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func connect(token string) {
	const url = "https://slack.com/api/rtm.connect?token=%s&presence_sub=1&batch_presence_aware=1"

	response, err := http.Get(fmt.Sprintf(url, token))

	if err != nil {
		// handle error
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	var v interface{}
	json.Unmarshal(body, &v)
	data := v.(map[string]interface{})
	// {"ok":true,"url":"wss:\/\/cerberus-xxxx.lb.slack-msgs.com\/websocket\/FI7ykBRQQ6b7ylhofO-HcuF0iP6ZmdrNCbM3Dlbdxnkco3TSWmGvYy5ZoF0Upibbps7zmynedGuVnkeuTE_wK3qKw9SwW20aG4q65LBEW0Y=","team":{"id":"T1CAKHC0N","name":"singmyr.io","domain":"singmyr"},"self":{"id":"UKWTYQEJV","name":"idle_rpg"}}
	for k, v := range data {
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
	}
	fmt.Println(data["url"])
}

func main() {
	connect("INSERT_TOKEN_HERE")
}
