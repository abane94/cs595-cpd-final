package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// groupme bot id used to get and post responses
const BotID = "c947489260df70ef469fc34ea8"

// a generic interface so different structs can be passed
// to the post function reguardless of their contents
// this might not be needed idk
type jsonBody interface {
}

// json format sent from groupMe for each msg in the group
type groupMe_message_post struct {
	Group_id string // `json:"group_id"`		// the names have to capitalized to be 'exported'
	Name     string // `json:"name"`
	Text     string // `json:"text"`
}

// what each message we send to group me will look like
type groupMe_message_send struct {
	bot_id string
	text   string
}

// the json respone from hitting the translating service
type translated_respone struct {
	Code int64
	Lang string
	Text []string
}

// the internal struct we can use to send data through the pipe
// not to be sent/received via json
type reply_data struct {
}

// sends post data to url with body converted to json
func post(url string, body jsonBody) {
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(body)
	//res, err := http.Post(url, "application/json; charset=utf-8", b)
	//if err != nil {
	//log some error stuff
	//}
}

// the thread the will synchronously get the translation and send it to the group
//func reply(pipe chan reply_data) {
//var data <- pipe
//for true {
// do stuff
// get data from translation endpoint
// send translation to groupMe room

//data <- pipe
//}
//}

// groupMe api of what is posted to our server
// {
// 	"attachments": [],
// 	"avatar_url": "https://i.groupme.com/123456789",
// 	"created_at": 1302623328,
// 	"group_id": "1234567890",
// 	"id": "1234567890",
// 	"name": "John",
// 	"sender_id": "12345",
// 	"sender_type": "user",
// 	"source_guid": "GUID",
// 	"system": false,
// 	"text": "Hello world ☃☃",
// 	"user_id": "1234567890"
//   }

func main() {
	// create the channel
	//pipe := make(chan reply_data, 1) // not sure what the second param does

	// start the reply thread
	//go reply(pipe) // run the reply goroutine

	// needed for heroku, gets the port
	port := os.Getenv("PORT")

	var groupMeMsg groupMe_message_send
	groupMeMsg.bot_id = BotID

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	// base get endpoint
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Team 2 says:<br/>Hello, you've requested tacos: %s\n", r.URL.Path)
	})

	// our message endpoint, where groupMe send the post
	http.HandleFunc("/group-me/msg", func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body) // collects the body
		if err != nil {
			log.Println(err)
			return
		}
		log.Println(string(body))

		var msg groupMe_message_post
		err = json.Unmarshal(body, &msg) // parses the json from the body
		if err != nil {
			log.Println(err)
			return
		}
		log.Println(msg)
		log.Println("LOG: msg.text :")
		log.Println(msg.Text)

		if strings.HasPrefix(msg.Text, "translate: ") {
			var msgToTranslate = msg.Text[11:len(msg.Text)]
			msgToTranslate = strings.Replace(msgToTranslate, " ", "%20", -1)
			log.Println(msgToTranslate)
			url := "https://translate.yandex.net/api/v1.5/tr.json/translate?" +
				"key=trnsl.1.1.20180404T152827Z.de1604f76f1d895c.8909d7acdac0907096f9a3cac7ecd33db02e0650&lang=en-de" +
				"&text=%20" + msgToTranslate
			client := http.Client{
				Timeout: time.Second * 2, // Maximum of 2 secs
			}
			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				log.Fatal(err)
				return
			}
			resp1, err1 := client.Do(req)
			if err1 != nil {
				log.Fatal(err1)
				return
			}
			log.Println("Sucessfully posted to yandex")
			body, readErr := ioutil.ReadAll(resp1.Body)
			if readErr != nil {
				log.Fatal(readErr)
			}
			log.Println("After reading in response from yandex")
			var jsonResponse translated_respone
			jsonErr := json.Unmarshal(body, &jsonResponse)
			if jsonErr != nil {
				log.Printf("error decoding yandex response: %v", err)
				if e, ok := err.(*json.SyntaxError); ok {
					log.Printf("yandex error at byte offset %d", e.Offset)
				}
				log.Printf("yandex response: %q", body)
				return
			}
			log.Println("json:")
			log.Println(jsonResponse)
			var translatedMsg = jsonResponse.Text
			log.Println(translatedMsg)
			var singleTranslatedMsg = strings.Replace(translatedMsg[0], " ", "%20", -1)
			resp2, err2 := http.Post("https://api.groupme.com/v3/bots/post?"+
				"bot_id="+BotID+
				"&text="+singleTranslatedMsg, "", nil)
			if err2 != nil {
				log.Fatal(err2)
			}
			defer resp2.Body.Close()
		}
	})

	http.ListenAndServe(":"+port, nil)
}
