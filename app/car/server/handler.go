package server

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"

	"github.com/shanghuiyang/rpi-devices/app/car/selfdriving"
	"github.com/shanghuiyang/rpi-devices/app/car/selfnav"
	"github.com/shanghuiyang/rpi-devices/app/car/selftracking"
	"github.com/shanghuiyang/rpi-devices/app/car/speechdriving"
	"github.com/shanghuiyang/rpi-devices/util"
	"github.com/shanghuiyang/rpi-devices/util/geo"

	"gocv.io/x/gocv"
)

const (
	ipPattern          = "((000.000.000.000))"
	selfDrivingState   = "((selfdriving-state))"
	selfTrackingState  = "((selftracking-state))"
	speechDrivingState = "((speechdriving-state))"
	volumePattern      = "((current-volume))"

	selfDrivingEnabled   = "((selfdriving-enabled))"
	selfTrackingEnabled  = "((selftracking-enabled))"
	speechDrivingEnabled = "((speechdriving-enabled))"

	logHandlerTag = "handler"
)

var (
	ip          string
	pageContext []byte
)

func (s *service) loadHomeHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[%v]load home page", logHandlerTag)
	rbuf := bytes.NewBuffer(pageContext)
	wbuf := bytes.NewBuffer([]byte{})
	volume, err := util.GetVolume()
	if err != nil {
		log.Printf("[%v]failed to get volume, error: %v", logHandlerTag, err)
		volume = 40
	}
	disabled := false
	selfDriving := selfdriving.Status()
	selfTracking := selftracking.Status()
	speechDriving := speechdriving.Status()
	if selfDriving || selfTracking || speechDriving {
		disabled = true
	}

	for {
		line, err := rbuf.ReadBytes('\n')
		if err == io.EOF {
			break
		}
		sline := string(line)

		sline = strings.Replace(sline, ipPattern, ip, 1)
		sline = strings.Replace(sline, volumePattern, fmt.Sprintf("%v", volume), 1)
		if selfTracking {
			sline = strings.Replace(sline, s.cfg.VideoHost, s.cfg.SelfTracking.VideoHost+"/video", 1)
		}

		if strings.Contains(sline, selfDrivingState) {
			state := "unchecked"
			if selfDriving {
				state = "checked"
			}
			sline = strings.Replace(sline, selfDrivingState, state, 1)

			able := "enabled"
			if state == "unchecked" && disabled {
				able = "disabled"
			}
			sline = strings.Replace(sline, selfDrivingEnabled, able, 1)
		}

		if strings.Contains(sline, selfTrackingState) {
			state := "unchecked"
			if selfTracking {
				state = "checked"
			}
			sline = strings.Replace(sline, selfTrackingState, state, 1)

			able := "enabled"
			if state == "unchecked" && disabled {
				able = "disabled"
			}
			sline = strings.Replace(sline, selfTrackingEnabled, able, 1)
		}

		if strings.Contains(sline, speechDrivingState) {
			state := "unchecked"
			if speechDriving {
				state = "checked"
			}
			sline = strings.Replace(sline, speechDrivingState, state, 1)

			able := "enabled"
			if state == "unchecked" && disabled {
				able = "disabled"
			}
			sline = strings.Replace(sline, speechDrivingEnabled, able, 1)
		}

		wbuf.Write([]byte(sline))
	}
	w.Write(wbuf.Bytes())
}

func (s *service) opHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[%v]op", logHandlerTag)
	vars := mux.Vars(r)
	v, ok := vars["op"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "invalid op: %v", vars["op"])
		return
	}
	op := Op(v)
	s.chOp <- op
}

func (s *service) setVolumeHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[%v]set volume", logHandlerTag)
	vars := mux.Vars(r)
	v, err := strconv.Atoi(vars["v"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "invalid volume: %v", vars["v"])
		return
	}
	if v < 0 || v > 100 { // volume should be 0~100%
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "invalid volume: %v", vars["v"])
		return
	}

	log.Printf("[%v]set volume: %v%%", logHandlerTag, v)
	if err := util.SetVolume(v); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "server internal error: %v", err)
		return
	}
	util.PlayWav("current-volume.wav")
}

func (s *service) selfDrivingOnHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[%v]self-driving on", logHandlerTag)
	if !s.cfg.SelfDriving.Enabled {
		log.Printf("[%v]self-driving was disabled", logHandlerTag)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("self-driving was disabled"))
		return
	}
	selfdriving.Start()
}

func (s *service) selfDrivingOffHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[%v]self-driving off", logHandlerTag)
	selfdriving.Stop()
}

func (s *service) selfTrackingOnHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[%v]self-tracking on", logHandlerTag)
	if !s.cfg.SelfTracking.Enabled {
		log.Printf("[%v]self-tracking was disabled", logHandlerTag)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("self-tracking was disabled"))
		return
	}

	if selftracking.Status() {
		log.Printf("[%v]self-tracking is running", logHandlerTag)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("self-tracking is running"))
		return
	}

	if err := util.StopMotion(); err != nil {
		log.Printf("[%v]failed to stop motion server", logHandlerTag)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to stop motion service"))
	}

	chImg := make(chan *gocv.Mat, 64)
	defer func() {
		close(chImg)
	}()
	go selftracking.Start(chImg)

	cam, err := gocv.OpenVideoCapture(0)
	if err != nil {
		log.Printf("[%v]failed to new camera", logHandlerTag)
		return
	}
	defer cam.Close()
	cam.Set(gocv.VideoCaptureFocus, cam.ToCodec("MJPG"))
	cam.Set(gocv.VideoCaptureFPS, 25)
	cam.Set(gocv.VideoCaptureFrameWidth, 640)
	cam.Set(gocv.VideoCaptureFrameHeight, 480)

	img := gocv.NewMat()
	defer img.Close()

	for selftracking.Status() {
		util.DelayMs(200)
		cam.Grab(10)
		if !cam.Read(&img) {
			log.Printf("[%v]failed to read img from camera", logHandlerTag)
			continue
		}
		im := img.Clone()
		chImg <- &im
	}

}

func (s *service) selfTrackingOffHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[%v]self-tracking off", logHandlerTag)
	selftracking.Stop()
	util.DelayMs(1000)
	if err := util.StartMotion(); err != nil {
		log.Printf("[%v]failed to start motion server", logHandlerTag)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to start motion service"))
		return
	}
}

func (s *service) speechDrivingOnHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[%v]speech-driving on", logHandlerTag)
	if !s.cfg.SpeechDriving.Enabled {
		log.Printf("[%v]speech-driving was disabled", logHandlerTag)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("speech-driving was disabled"))
		return
	}
	s.ledBlinked = false
	speechdriving.Start()
}

func (s *service) speechDrivingOffHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[%v]speech-driving off", logHandlerTag)
	speechdriving.Stop()
	s.ledBlinked = true
	go s.blink()
}

func (s *service) selfNavOnHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[%v]self-nav on", logHandlerTag)
	if !s.cfg.SelfNav.Enabled {
		log.Printf("[%v]self-nav was disabled", logHandlerTag)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("self-nav was disabled"))
		return
	}
	vars := mux.Vars(r)
	lat, err := strconv.ParseFloat(vars["lat"], 64)
	if err != nil {
		log.Printf("[%v]invalid lat: %v", logHandlerTag, vars["lat"])
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "invalid lat: %v", vars["lat"])
		return
	}
	lon, err := strconv.ParseFloat(vars["lon"], 64)
	if err != nil {
		log.Printf("[%v]invalid lon: %v", logHandlerTag, vars["lon"])
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "invalid lat: %v", vars["lon"])
		return
	}

	if lat < -90 || lat > 90 { // volume should be 0~100%
		log.Printf("[%v]invalid lat: %v", logHandlerTag, lat)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "invalid lat: %v", lat)
		return
	}

	if lon < -180 || lon > 180 { // volume should be 0~100%
		log.Printf("[%v]invalid lon: %v", logHandlerTag, lon)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "invalid lon: %v", lon)
		return
	}

	dest := &geo.Point{
		Lat: lat,
		Lon: lon,
	}
	selfnav.Start(dest)
}

func (s *service) selfNavOffHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[%v]self-nav off", logHandlerTag)
	selfnav.Stop()
}

func (s *service) musicOnHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[car]music on")
	util.PlayMp3("./music/*.mp3")
}

func (s *service) musicOffHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[car]music off")
	util.StopMp3()
}
