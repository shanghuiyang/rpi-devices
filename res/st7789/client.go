package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func main() {
	img, err := ioutil.ReadFile("test.png")
	if err != nil {
		fmt.Printf("open map.png error: %v", err)
		return
	}
	buf := bytes.NewBuffer(img)
	req, err := http.NewRequest("POST", "http://localhost:8080/display", buf)
	if err != nil {
		fmt.Printf("failed to new http request, error: %v", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{
		Timeout: 2 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("failed to send http request, error: %v", err)
		return
	}
	defer resp.Body.Close()
}
