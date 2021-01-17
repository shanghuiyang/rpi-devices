package car

const (
	chSize        = 8
	letMeThinkWav = "let_me_think.wav"
	thisIsXWav    = "this_is_x.wav"
	iDontKnowWav  = "i_dont_know.wav"
	errorWav      = "error.wav"
)

const (
	baiduSpeechAppKey            = "your_speech_app_key"
	baiduSpeechSecretKey         = "your_speech_secret_key"
	baiduImgRecognitionAppKey    = "your_image_recognition_app_key"
	baiduImgRecognitionSecretKey = "your_image_recognition_secrect_key"
)

const (
	forward          Op = "forward"
	backward         Op = "backward"
	left             Op = "left"
	right            Op = "right"
	stop             Op = "stop"
	pause            Op = "pause"
	turn             Op = "turn"
	scan             Op = "scan"
	roll             Op = "roll"
	beep             Op = "beep"
	blink            Op = "blink"
	servoleft        Op = "servoleft"
	servoright       Op = "servoright"
	servoahead       Op = "servoahead"
	lighton          Op = "lighton"
	lightoff         Op = "lightoff"
	musicon          Op = "musicon"
	musicoff         Op = "musicoff"
	selfdrivingon    Op = "selfdrivingon"
	selfdrivingoff   Op = "selfdrivingoff"
	selftrackingon   Op = "selftrackingon"
	selftrackingoff  Op = "selftrackingoff"
	speechdrivingon  Op = "speechdrivingon"
	speechdrivingoff Op = "speechdrivingoff"
	selfnavon        Op = "selfnavon"
	selfnavoff       Op = "selfnavoff"
)

var (
	scanningAngles = []int{-90, -75, -60, -45, -30, -15, 0, 15, 30, 45, 60, 75, 90}
	aheadAngles    = []int{0, -15, 0, 15}
)

const (
	// the hsv of a tennis
	lh float64 = 33
	ls float64 = 108
	lv float64 = 138
	hh float64 = 61
	hs float64 = 255
	hv float64 = 255
)

type (
	// Op ...
	Op string
	// Option ...
	Option func(c *Car)
)
