![](img/go-devices.png)

# rpi-devices 
[![ci](https://github.com/shanghuiyang/rpi-devices/actions/workflows/ci.yml/badge.svg)](https://github.com/shanghuiyang/rpi-devices/actions/workflows/ci.yml)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](https://github.com/shanghuiyang/rpi-devices/blob/master/LICENSE)

rpi-devices implements drivers for various kinds of sensors or devices based on raspberry pi in pure golang.

## Usage
```go
package main

import (
	"time"

	"github.com/shanghuiyang/rpi-devices/dev"
)

const (
	pin = 26
)

func main() {
	led := dev.NewLedImp(pin)

	led.On()
	time.Sleep(3 * time.Second)
	led.Off()
}
```

## Currently Implemented Drivers

|Sensors|Image|Description|Example|App|
|-------|-----|-----|-------|---|
|ADS1015|![](img/ads1015.jpg)|Analog-to-digital converter|N/A|[joystick](/app/joystick)|
|Button|![](img/button.jpg)|Button module|[example](/example/button/main.go)|[vedio-monitor](/app/vmonitor)|
|Buzzer|![](img/buzzer.jpg)|Buzzer module|N/A|[car](/app/car), [door-dog](/app/doordog)|
|BYJ2848|![](img/step-motor.jpg)|Step motor|[example](/example/byj2848/main.go)|N/A|
|Collision Switch|![](img/collision-switch.jpg)|A switch for deteching collision|[example](/example/collision_switch/main.go)|[car](/app/car)|
|DHT11|![](img/dht11.jpg)|Temperature & Humidity sensor|[example](/example/dht11/main.go)|[home-asst](/app/homeasst)|
|Digital Led Display|![](img/digital-led-display.jpg)|led digital module|[example](/example/digital_led_display/main.go)|[auto-air](/app/autoair)|
|DS18B20|![](img/temp.jpg)|Temperature sensor|[example](/example/temperature/main.go)|[auto-fan](/app/autofan)|
|Encoder|![](img/encoder.jpg)|Encoder sensor|[example](/example/encoder/main.go)|N/A|
|GPS|![](img/gps.jpg)|Location sensor|[example](/example/gps/main.go)|[gps-tracker](/app/gpstracker)|
|GY-25|![](img/gy25.jpg)|Angle sensor|[example](/example/gy25/main.go)|[car](/app/car)|
|HC-SR04|![](img/hc-sr04.jpg)|Ultrasonic distance meter|[example](/example/hcsr04/main.go)|[auto-light](/app/autolight), [doordog](/app/doordog)|
|Humidity Detector|![](img/humidity-detector.jpg)|Soil humidity detector|[example](/example/humidity_detector/main.go)|N/A|
|Infrared Encoder/Decoder|![](img/ir-encoder-decoder.jpg)|Infrared encoder/decoder|[example](/example/ir_coder/main.go)|N/A|
|Infrared|![](img/infared.jpg)|Infrared sensor|[example](/example/ir_detector/main.go)|N/A|
|Joystick|![](img/joystick.jpg)|XY Dual Axis Joystick|[example](/example/joystick/main.go)|[car](/app/car)|
|L298N|![](img/l298n.jpg)|motor driver|N/A|[car](/app/car)|
|LC12S|![](img/lc12s.jpg)|2.4g wireless module|[example](/example/lc12s/main.go)|[car](/app/car)|
|LCD 1602A Display|![](img/lcd1602a.jpg)|lcd display module|[example](/example/lcd_display/main.go)|[home-asst](/app/homeasst)
|Led|![](img/led.jpg)|Led light|[example](/example/led/main.go)|[car](/app/car), [vedio-monitor](/app/vmonitor)|
|MPU6050|![](img/mpu6050.jpg)|6-axis motion sensor|[example](/example/mpu6050/main.go)|N/A|
|Oled|![](img/oled.jpg)|Oled display module|[example](/example/oled_display/main.go)|[home-asst](/app/homeasst)|
|PCF8591|![](img/pcf8591.jpg)|Analog-to-digital converter|N/A|N/A|
|PMS7003|![](img/pms7003.jpg)|Air quality sensor|[example](/example/air/main.go)|[auto-air](/app/autoair)|
|Relay|![](img/relay.jpg)|Relay module|[example](/example/relay/main.go)|[auto-fan](/app/autofan)|
|RX480E-4|![](img/rx480e4.jpg)|433MHz Wireless RF Receiver|[example](/example/rx480e4/main.go)|[remote-light](/app/rlight)|
|SG90|![](img/sg90.jpg)|Servo motor|[example](/example/sg90/main.go)|[auto-air](/app/autoair), [car](/app/car), [vedio-monitor](/app/vmonitor)|
|SW-420|![](img/sw-420.jpg)|Shaking sensor|[example](/example/sw420/main.go)|[auto-air-out](/app/autoairout)|
|US-100|![](img/us-100.jpg)|ultrasonic distance meter|[example](/example/us100/main.go)|[car](/app/car)|
|Voice|![](img/voice.jpg)|Voice sensor|N/A|N/A|
|ZE08-CH2O|![](img/ze08-ch2o.jpg)|CH2O sensor|[example](/example/ze08ch2o/main.go)|[ch2o-monitor](/app/ch2omonitor)|

## Apps
Using the driver programs, I built several applications. The most complex app is the [smart car](/app/car), more than 10 sensers were used to build the car. I highlight few funny apps here, please go to [app](/app) for all apps I developed. You can learn how to use the drivers from my apps.
### [Self-Dirving Car](/app/car)
play the video on [youtube](https://www.youtube.com/watch?v=RNqe4byzXmw).

<img src="img/car.gif" width=60% height=60% />

### [Video Monitor](/app/vmonitor)
<img src="img/vmonitor.gif" width=60% height=60% />

### [Auto-Air](/app/autoair)
<img src="img/auto-air.gif" width=60% height=60% />

### [GPS-Tracker](/app/gpstracker)
<img src="img/gpstracker-v2.gif" width=60% height=60% />

### [Home-Asst](/app/homeasst)
<img src="img/homeasst.gif" width=40% height=40% />
