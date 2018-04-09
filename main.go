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
)


func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Team 2 says:<br/>Hello, you've requested: %s\n", r.URL.Path)
	})

	http.ListenAndServe(":" + port, nil)
}