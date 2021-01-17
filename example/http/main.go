/*
Auto-Air opens the air-cleaner automatically when the pm2.5 >= 130.
*/

package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", server)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err.Error())
	}
}

func server(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello world!"))
}
