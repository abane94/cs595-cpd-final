// package main

// import (
// 	"log"
// 	"net/http"
// 	"os"

// 	"github.com/gin-gonic/gin"
// 	_ "github.com/heroku/x/hmetrics/onload"
// )

// func main() {
// 	port := os.Getenv("PORT")

// 	if port == "" {
// 		log.Fatal("$PORT must be set")
// 	}

// 	router := gin.New()
// 	router.Use(gin.Logger())
// 	router.LoadHTMLGlob("templates/*.tmpl.html")
// 	router.Static("/static", "static")

// 	router.GET("/", func(c *gin.Context) {
// 		c.HTML(http.StatusOK, "index.tmpl.html", nil)
// 	})

// 	router.Run(":" + port)
// }

package main

import (
	"fmt"
	"net/http"
	"log"
	"os"
	"encoding/json"
	"io/ioutil"
)

type message struct {
	group_id string `json:"group_id"`
	name string `json:"name"`
	text string `json:"text"`
}

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
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Team 2 says:<br/>Hello, you've requested: %s\n", r.URL.Path)
	})

	http.HandleFunc("/group-me/msg", func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
			return
		}
		log.Println(string(body))
		
		msg := message{}
		err = json.Unmarshal(body, &msg)
		if err != nil {
			// fmt.Println(err)
			log.Println(err)
			return
		}
		log.Println(r.Body)
		// fmt.Fprintf(w, "Team 2 says:\nHello, you've requested: %s\n", r.URL.Path)
		fmt.Fprintf(w, msg.text)
	})

	http.ListenAndServe(":" + port, nil)
}