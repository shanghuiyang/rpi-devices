package gocv

import (
	"errors"
)

// VideoCaptureProperties are the properties used for VideoCapture operations.
type VideoCaptureProperties int

const (
	// VideoCapturePosMsec contains current position of the
	// video file in milliseconds.
	VideoCapturePosMsec VideoCaptureProperties = 0

	// VideoCapturePosFrames 0-based index of the frame to be
	// decoded/captured next.
	VideoCapturePosFrames VideoCaptureProperties = 1

	// VideoCapturePosAVIRatio relative position of the video file:
	// 0=start of the film, 1=end of the film.
	VideoCapturePosAVIRatio VideoCaptureProperties = 2

	// VideoCaptureFrameWidth is width of the frames in the video stream.
	VideoCaptureFrameWidth VideoCaptureProperties = 3

	// VideoCaptureFrameHeight controls height of frames in the video stream.
	VideoCaptureFrameHeight VideoCaptureProperties = 4

	// VideoCaptureFPS controls capture frame rate.
	VideoCaptureFPS VideoCaptureProperties = 5

	// VideoCaptureFOURCC contains the 4-character code of codec.
	// see VideoWriter::fourcc for details.
	VideoCaptureFOURCC VideoCaptureProperties = 6

	// VideoCaptureFrameCount contains number of frames in the video file.
	VideoCaptureFrameCount VideoCaptureProperties = 7

	// VideoCaptureFormat format of the Mat objects returned by
	// VideoCapture::retrieve().
	VideoCaptureFormat VideoCaptureProperties = 8

	// VideoCaptureMode contains backend-specific value indicating
	// the current capture mode.
	VideoCaptureMode VideoCaptureProperties = 9

	// VideoCaptureBrightness is brightness of the image
	// (only for those cameras that support).
	VideoCaptureBrightness VideoCaptureProperties = 10

	// VideoCaptureContrast is contrast of the image
	// (only for cameras that support it).
	VideoCaptureContrast VideoCaptureProperties = 11

	// VideoCaptureSaturation saturation of the image
	// (only for cameras that support).
	VideoCaptureSaturation VideoCaptureProperties = 12

	// VideoCaptureHue hue of the image (only for cameras that support).
	VideoCaptureHue VideoCaptureProperties = 13

	// VideoCaptureGain is the gain of the capture image.
	// (only for those cameras that support).
	VideoCaptureGain VideoCaptureProperties = 14

	// VideoCaptureExposure is the exposure of the capture image.
	// (only for those cameras that support).
	VideoCaptureExposure VideoCaptureProperties = 15

	// VideoCaptureConvertRGB is a boolean flags indicating whether
	// images should be converted to RGB.
	VideoCaptureConvertRGB VideoCaptureProperties = 16

	// VideoCaptureWhiteBalanceBlueU is currently unsupported.
	VideoCaptureWhiteBalanceBlueU VideoCaptureProperties = 17

	// VideoCaptureRectification is the rectification flag for stereo cameras.
	// Note: only supported by DC1394 v 2.x backend currently.
	VideoCaptureRectification VideoCaptureProperties = 18

	// VideoCaptureMonochrome indicates whether images should be
	// converted to monochrome.
	VideoCaptureMonochrome VideoCaptureProperties = 19

	// VideoCaptureSharpness controls image capture sharpness.
	VideoCaptureSharpness VideoCaptureProperties = 20

	// VideoCaptureAutoExposure controls the DC1394 exposure control
	// done by camera, user can adjust reference level using this feature.
	VideoCaptureAutoExposure VideoCaptureProperties = 21

	// VideoCaptureGamma controls video capture gamma.
	VideoCaptureGamma VideoCaptureProperties = 22

	// VideoCaptureTemperature controls video capture temperature.
	VideoCaptureTemperature VideoCaptureProperties = 23

	// VideoCaptureTrigger controls video capture trigger.
	VideoCaptureTrigger VideoCaptureProperties = 24

	// VideoCaptureTriggerDelay controls video capture trigger delay.
	VideoCaptureTriggerDelay VideoCaptureProperties = 25

	// VideoCaptureWhiteBalanceRedV controls video capture setting for
	// white balance.
	VideoCaptureWhiteBalanceRedV VideoCaptureProperties = 26

	// VideoCaptureZoom controls video capture zoom.
	VideoCaptureZoom VideoCaptureProperties = 27

	// VideoCaptureFocus controls video capture focus.
	VideoCaptureFocus VideoCaptureProperties = 28

	// VideoCaptureGUID controls video capture GUID.
	VideoCaptureGUID VideoCaptureProperties = 29

	// VideoCaptureISOSpeed controls video capture ISO speed.
	VideoCaptureISOSpeed VideoCaptureProperties = 30

	// VideoCaptureBacklight controls video capture backlight.
	VideoCaptureBacklight VideoCaptureProperties = 32

	// VideoCapturePan controls video capture pan.
	VideoCapturePan VideoCaptureProperties = 33

	// VideoCaptureTilt controls video capture tilt.
	VideoCaptureTilt VideoCaptureProperties = 34

	// VideoCaptureRoll controls video capture roll.
	VideoCaptureRoll VideoCaptureProperties = 35

	// VideoCaptureIris controls video capture iris.
	VideoCaptureIris VideoCaptureProperties = 36

	// VideoCaptureSettings is the pop up video/camera filter dialog. Note:
	// only supported by DSHOW backend currently. The property value is ignored.
	VideoCaptureSettings VideoCaptureProperties = 37

	// VideoCaptureBufferSize controls video capture buffer size.
	VideoCaptureBufferSize VideoCaptureProperties = 38

	// VideoCaptureAutoFocus controls video capture auto focus..
	VideoCaptureAutoFocus VideoCaptureProperties = 39

	// VideoCaptureSarNumerator controls the sample aspect ratio: num/den (num)
	VideoCaptureSarNumerator VideoCaptureProperties = 40

	// VideoCaptureSarDenominator controls the sample aspect ratio: num/den (den)
	VideoCaptureSarDenominator VideoCaptureProperties = 41

	// VideoCaptureBackend is the current api backend (VideoCaptureAPI). Read-only property.
	VideoCaptureBackend VideoCaptureProperties = 42

	// VideoCaptureChannel controls the video input or channel number (only for those cameras that support).
	VideoCaptureChannel VideoCaptureProperties = 43

	// VideoCaptureAutoWB controls the auto white-balance.
	VideoCaptureAutoWB VideoCaptureProperties = 44

	// VideoCaptureWBTemperature controls the white-balance color temperature
	VideoCaptureWBTemperature VideoCaptureProperties = 45

	// VideoCaptureCodecPixelFormat shows the the codec's pixel format (4-character code). Read-only property.
	// Subset of AV_PIX_FMT_* or -1 if unknown.
	VideoCaptureCodecPixelFormat VideoCaptureProperties = 46

	// VideoCaptureBitrate displays the video bitrate in kbits/s. Read-only property.
	VideoCaptureBitrate VideoCaptureProperties = 47
)

type VideoCapture struct{}

func OpenVideoCapture(v interface{}) (*VideoCapture, error) {
	return nil, errors.New("not implement")
}

func (v *VideoCapture) Grab(skip int) {}

func (v *VideoCapture) Read(m *Mat) bool {
	return false
}
func (v *VideoCapture) Set(prop VideoCaptureProperties, param float64) {}

func (v *VideoCapture) ToCodec(codec string) float64 {
	return 0
}

func (v *VideoCapture) Close() error {
	return errors.New("not implement")
}
