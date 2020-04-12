/*
Auto-Air opens the air-cleaner automatically when the pm2.5 >= 130.
*/

package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto"
)

type alarm string

const (
	alarmClock       alarm = "AlarmClock"
	policeSiren      alarm = "PoliceSiren"
	railroadCrossing alarm = "RailroadCrossing"
)

var sounds = map[alarm]string{
	alarmClock:       "sounds/alarm_clock.mp3",
	policeSiren:      "/Users/shanghui.yang/Downloads/classic.mp3",
	railroadCrossing: "sounds/RailroadCrossing.mp3",
}

type server struct {
	ch       chan alarm
	alarming bool
}

func main() {
	svr := newServer()
	svr.start()
}

func newServer() *server {
	return &server{
		ch: make(chan alarm, 256),
	}
}

func (s *server) start() {
	go s.run()

	log.Printf("server started")
	http.HandleFunc("/", s.handler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err.Error())
	}
}

func (s *server) run() {
	for a := range s.ch {
		if s.alarming {
			// log.Printf("skip a request")
			continue
		}
		s.alarming = true
		go s.play(a)
	}
}

func (s *server) handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		w.Write([]byte("hello"))
	case "POST":
		fv := r.FormValue("alarm")
		// log.Printf("request: %v", fv)
		a := alarm(fv)
		s.handle(a)
	}
}

func (s *server) handle(a alarm) {
	s.ch <- a
}

func (s *server) play(a alarm) {
	log.Printf("start alarm")
	defer func() {
		time.Sleep(1 * time.Second)
		s.alarming = false
		log.Printf("finish alarm")
	}()

	mp3f := sounds[alarmClock]
	if f, ok := sounds[a]; ok {
		mp3f = f
	}
	if err := s.playmp3(mp3f); err != nil {
		log.Fatal(err)
	}
}

func (s *server) playmp3(mp3f string) error {
	f, err := os.Open(mp3f)
	if err != nil {
		return err
	}
	defer f.Close()

	d, err := mp3.NewDecoder(f)
	if err != nil {
		return err
	}

	c, err := oto.NewContext(d.SampleRate(), 2, 2, 8192)
	if err != nil {
		return err
	}
	defer c.Close()

	p := c.NewPlayer()
	defer p.Close()

	if _, err := io.Copy(p, d); err != nil {
		return err
	}
	return nil
}
