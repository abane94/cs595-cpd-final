
package main

import (
	"fmt"
	"net/http"
	"log"
	"os"
	"encoding/json"
	"io/ioutil"
)

const BOT_ID="" // need to get once its registered with groupMe

// a generic interface so different structs can be passed
// to the post function reguardless of their contents
// this might not be needed idk
type jsonBody interface{

}

// json format sent from groupMe for each msg in the group
type groupMe_message_post struct {
	Group_id string// `json:"group_id"`		// the names have to capitalized to be 'exported'
	Name string// `json:"name"`
	Text string// `json:"text"`
}

// what each message we send to group me will look like
type groupMe_message_send struct {
	// todo
}

// the json respone from hitting the translating service
type translated_respone struct {
	Status int64
	Lang string
	Text string[]
}

// the internal struct we can use to send data through the pipe
// not to be sent/received via json
type reply_data struct {

}


// sends post data to url with body converted to json
func post(url string, body jsonBody) {
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(body)
	res, err := http.Post(url, "application/json; charset=utf-8", b)
	if err != nil {
		//log some error stuff
	}
}

// the thread the will synchronously get the translation and send it to the group
func reply(pipe chan reply_data) {
	data <- pipe
	while (data.end != nil) {
		// do stuff
			// get data from translation endpoint
			// send translation to groupMe room

		data <- pipe
	}
}

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
	pipe := make(chan reply_data, 1) // not sure what the second param does

	// start the reply thread
	go reply(pipe) // run the reply goroutine


	// needed for heroku, gets the port
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}


	// base get endpoint
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Team 2 says:<br/>Hello, you've requested: %s\n", r.URL.Path)
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
			// fmt.Println(err)
			log.Println(err)
			return
		}
		// fmt.Fprintf(w, "Team 2 says:\nHello, you've requested: %s\n", r.URL.Path)
		log.Println(msg)
		log.Println("LOG: msg.text :")
		log.Println(msg.Text)

		// sends response, not sure how groupMe will handle this
		fmt.Fprintf(w, msg.Text)
	})

	http.ListenAndServe(":" + port, nil)
}