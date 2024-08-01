![](img/go-devices.png)

## rpi-devices 
[![ci](https://github.com/shanghuiyang/rpi-devices/actions/workflows/ci.yml/badge.svg)](https://github.com/shanghuiyang/rpi-devices/actions/workflows/ci.yml)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](https://github.com/shanghuiyang/rpi-devices/blob/master/LICENSE)
[![GoDoc](https://godoc.org/github.com/shanghuiyang/rpi-devices?status.svg)](https://godoc.org/github.com/shanghuiyang/rpi-devices)

rpi-devices implements drivers for various kinds of sensors or devices based on [raspberry pi](https://www.raspberrypi.org/) in pure Golang.

### Usage
```go
package main

import (
	"time"

	"github.com/shanghuiyang/rpi-devices/dev"
)

const pin = 26

func main() {
	led := dev.NewLedImp(pin)
	for {
		led.On()
		time.Sleep(1 * time.Second)
		led.Off()
		time.Sleep(1 * time.Second)
	}
}
```

### Currently Implemented Drivers

|Sensors|Image|Description|Example|Projects|
|-------|-----|-----|-------|---|
|ADS1015|![](img/ads1015.jpg)|Analog-to-digital converter|N/A|[joystick](https://github.com/shanghuiyang/rpi-projects/tree/main/projects/joystick)|
|Button|![](img/button.jpg)|Button module|[example](/example/button/main.go)|[vedio-monitor](https://github.com/shanghuiyang/rpi-projects/tree/main/projects/vmonitor)|
|Buzzer|![](img/buzzer.jpg)|Buzzer module|N/A|[car](https://github.com/shanghuiyang/rpi-projects/tree/main/projects/car), [door-dog](https://github.com/shanghuiyang/rpi-projects/tree/main/projects/doordog)|
|BYJ2848|![](img/step-motor.jpg)|Step motor|[example](/example/byj2848/main.go)|N/A|
|Collision Switch|![](img/collision-switch.jpg)|A switch for deteching collision|[example](/example/collision_switch/main.go)|[car](https://github.com/shanghuiyang/rpi-projects/tree/main/projects/car)|
|DHT11|![](img/dht11.jpg)|Temperature & Humidity sensor|[example](/example/dht11/main.go)|[home-asst](https://github.com/shanghuiyang/rpi-projects/tree/main/projects/homeasst)|
|Display Digital Led TM1637 |![](img/digital-led-display.jpg)|Digital led module|[example](/example/display_led_tm1637/main.go)|[auto-air](https://github.com/shanghuiyang/rpi-projects/tree/main/projects/autoair)|
|Display LCD|![](img/lcd1602a.jpg)|LCD display module|[example](/example/display_lcd/main.go)|[home-asst](https://github.com/shanghuiyang/rpi-projects/tree/main/projects/homeasst)
|Display SSD1360|![](img/oled.jpg)|Oled display module|[example](/example/display_oled_ssd1306/main.go)|[home-asst](https://github.com/shanghuiyang/rpi-projects/tree/main/projects/homeasst)|
|Display ST7899|![](img/tft_st7899.jpg)|TFT LCD display module|[example](/example/display_oled_ssd1306/main.go)|[gps-tracker](https://github.com/shanghuiyang/rpi-projects/tree/main/projects/gpstracker)|
|DS18B20|![](img/temp.jpg)|Temperature sensor|[example](/example/temperature/main.go)|[auto-fan](https://github.com/shanghuiyang/rpi-projects/tree/main/projects/autofan)|
|Encoder|![](img/encoder.jpg)|Encoder sensor|[example](/example/encoder/main.go)|N/A|
|GPS NEO-6M|![](img/gps-neo6m.jpg)|Location sensor|[example](/example/gps/main.go)|[car](https://github.com/shanghuiyang/rpi-projects/tree/main/projects/car)|
|GPS HT1818|![](img/gps-ht1818.jpg)|Location sensor|N/A|[gps-tracker](https://github.com/shanghuiyang/rpi-projects/tree/main/projects/gpstracker)|
|GY-25|![](img/gy25.jpg)|Angle sensor|[example](/example/gy25/main.go)|[car](https://github.com/shanghuiyang/rpi-projects/tree/main/projects/car)|
|HC-SR04|![](img/hc-sr04.jpg)|Ultrasonic distance meter|[example](/example/hcsr04/main.go)|[auto-light](https://github.com/shanghuiyang/rpi-projects/tree/main/projects/autolight), [doordog](https://github.com/shanghuiyang/rpi-projects/tree/main/projects/doordog)|
|HDC1080|![](img/hdc1080.jpg)|Thermohygrometer sensor|[example](/example/hdc1080/main.go)|[home-asst](https://github.com/shanghuiyang/rpi-projects/tree/main/projects/homeasst)|
|Humidity Detector|![](img/humidity-detector.jpg)|Soil humidity detector|[example](/example/humidity_detector/main.go)|N/A|
|Infrared Encoder/Decoder|![](img/ir-encoder-decoder.jpg)|Infrared encoder/decoder|[example](/example/ir_coder/main.go)|N/A|
|Infrared|![](img/infared.jpg)|Infrared sensor|[example](/example/ir_detector/main.go)|N/A|
|Joystick|![](img/joystick.jpg)|XY Dual Axis Joystick|[example](/example/joystick/main.go)|[car](https://github.com/shanghuiyang/rpi-projects/tree/main/projects/car)|
|L298N|![](img/l298n.jpg)|Motor driver|N/A|[car](https://github.com/shanghuiyang/rpi-projects/tree/main/projects/car)|
|LC12S|![](img/lc12s.jpg)|2.4g wireless module|[example](/example/lc12s/main.go)|[car](https://github.com/shanghuiyang/rpi-projects/tree/main/projects/car)|
|Led|![](img/led.jpg)|Led light|[example](/example/led/main.go)|[car](https://github.com/shanghuiyang/rpi-projects/tree/main/projects/car), [vedio-monitor](https://github.com/shanghuiyang/rpi-projects/tree/main/projects/vmonitor)|
|MPU6050|![](img/mpu6050.jpg)|6-axis motion sensor|[example](/example/mpu6050/main.go)|N/A|
|PCF8591|![](img/pcf8591.jpg)|Analog-to-digital converter|N/A|N/A|
|PMS7003|![](img/pms7003.jpg)|Air quality sensor|[example](/example/air/main.go)|[auto-air](https://github.com/shanghuiyang/rpi-projects/tree/main/projects/autoair)|
|Relay|![](img/relay.jpg)|Relay module|[example](/example/relay/main.go)|[auto-fan](https://github.com/shanghuiyang/rpi-projects/tree/main/projects/autofan)|
|RX480E-4|![](img/rx480e4.jpg)|433MHz Wireless RF Receiver|[example](/example/rx480e4/main.go)|[remote-light](https://github.com/shanghuiyang/rpi-projects/tree/main/projects/rlight)|
|SG90|![](img/sg90.jpg)|Servo motor|[example](/example/sg90/main.go)|[auto-air](https://github.com/shanghuiyang/rpi-projects/tree/main/projects/autoair), [car](https://github.com/shanghuiyang/rpi-projects/tree/main/projects/car), [vedio-monitor](https://github.com/shanghuiyang/rpi-projects/tree/main/projects/vmonitor)|
|SW-420|![](img/sw-420.jpg)|Shaking sensor|[example](/example/sw420/main.go)|[auto-air-out](https://github.com/shanghuiyang/rpi-projects/tree/main/projects/autoairout)|
|US-100|![](img/us-100.jpg)|Ultrasonic distance meter|[example](/example/us100/main.go)|[car](https://github.com/shanghuiyang/rpi-projects/tree/main/projects/car)|
|Voice Detector|![](img/voice.jpg)|Voice detector|N/A|N/A|
|Water Flow Sensor|![](img/water_flow_sensor.jpg)|Water flow sensor|[example](/example/water_flow_sensor/main.go)|N/A|
|ZE08-CH2O|![](img/ze08-ch2o.jpg)|CH2O sensor|[example](/example/ze08ch2o/main.go)|[ch2o-monitor](https://github.com/shanghuiyang/rpi-projects/tree/main/projects/ch2omonitor)|
|ZP16|![](img/zp16.jpg)|Gas detector|[example](/example/zp16/main.go)|[home-asst](https://github.com/shanghuiyang/rpi-projects/tree/main/projects/homeasst)|

### Projects
See my another repo [rpi-projects](https://github.com/shanghuiyang/rpi-projects) for all projects that I developed them using this libaray.
