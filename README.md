![](img/go-devices.png)

# rpi-devices 
[![Build Status](https://travis-ci.org/shanghuiyang/rpi-devices.svg?branch=master)](https://travis-ci.org/shanghuiyang/rpi-devices)

rpi-devices let you drive the devices and sensors using a raspberry pi in pure golang.
The following devices & sensors had been implemented in the current version, and I will keep implementing new devices and sensors.


|Sensors|Image|Description|Example|App|
|-------|-----|-----|-------|---|
|Button|![](img/button.jpg)|Button module|[example](/example/button/button.go)|[vedio-monitor](/app/vmonitor)|
|Buzzer|![](img/buzzer.jpg)|Buzzer module|N/A|[car](/app/car), [door-dog](/app/doordog)|
|Collision Switch|![](img/collision-switch.jpg)|A switch for deteching collision|[example](/example/collisionswitch/collisionswitch.go)|[car](/app/car)|
|DHT11|![](img/dht11.jpg)|Temperature & Humidity sensor|[example](/example/dht11/dht11.go)|[home-assit](/app/homeassit)|
|DS18B20|![](img/temp.jpg)|Temperature sensor|[example](/example/temperature/temperature.go)|[auto-fan](/app/autofan)|
|Encoder|![](img/encoder.jpg)|Encoder sensor|[example](/example/encoder/encoder.go)|[car](/app/car)|
|GPS|![](img/gps.jpg))|location sensor|[example](/example/gps/gps.go)|[gps-tracker](/app/gpstracker)|
|HC-SR04|![](img/hc-sr04.jpg)|ultrasonic distance meter|[example](/example/hcsr04/hcsr04.go)|[auto-light](/app/autolight), [doordog](/app/doordog)|
|Infrared|![](img/infared.jpg)|Infrared sensor|[example](/example/infrared/infrared.go)|N/A|
|L298N|![](img/l298n.jpg)|motor driver|N/A|[car](/app/car)|
|Led|![](img/led.jpg)|Led light|[example](/example/led/led.go)|[car](/app/car), [vedio-monitor](/app/vmonitor)|
|Led Display|![](img/digital-led-display.jpg)|led digital module|[example](/example/leddisplay/leddisplay.go)|[auto-air](/app/autoair)|
|Oled|![](img/oled.jpg)|Oled display module|[example](/example/oled/oled.go)|[home-assit](/app/homeassit)|
|PMS7003|![](img/pms7003.jpg)|Air quality sensor|[example](/example/air/air.go)|[auto-air](/app/autoair)|
|Relay|![](img/relay.jpg)|Relay module|[example](/example/relay/relay.go)|[auto-fan](/app/autofan)|
|RX480E-4|![](img/rx480e4.jpg)|Wireless remote control|[example](/example/rx480e4/rx480e4.go)|[remote-light](/app/rlight)|
|SG90|![](img/sg90.jpg)|Servo motor|[example](/example/sg90/sg90.go)|[auto-air](/app/autoair), [car](/app/car), [vedio-monitor](/app/vmonitor)|
|Step Motor|![](img/step-motor.jpg)|Step motor|[example](/example/stepmotor/stepmotor.go)|N/A|
|SW-420|![](img/sw-420.jpg)|Shaking sensor|[example](/example/sw420/sw420.go)|[auto-air-out](/app/autoairout)|
|US-100|![](img/us-100.jpg)|ultrasonic distance meter|[example](/example/us100/us100.go)|[car](/app/car)|
|Voice|![](img/voice.jpg)|Voice sensor|N/A|N/A|
|ZE08-CH2O|![](img/ze08-ch2o.jpg)|CH2O sensor|[example](/example/ch20/ch2o.go)|[ch2o-monitor](/app/ch2omonitor)|


## Compile & Deploy

It is very easy to cross-compile and deploy for golang. It is an example that compiles the binary for raspberry pi on MacOS.
```shell
$ CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=6 go build -o app.pi main.go
````

If you aren't sure the cpu info of your raspberry pi, check it out by,
```shell
$ lscpu
# those are the cpu info of my raspberry pi 2.
# ------------------------------------------------------------
# Architecture:        armv6l
# Byte Order:          Little Endian
# CPU(s):              1
# On-line CPU(s) list: 0
# Thread(s) per core:  1
# Core(s) per socket:  1
# Socket(s):           1
# Vendor ID:           ARM
# Model:               7
# Model name:          ARM1176
# Stepping:            r0p7
# CPU max MHz:         700.0000
# CPU min MHz:         700.0000
# BogoMIPS:            697.95
# Flags:               half thumb fastmult vfp edsp java tls
# ------------------------------------------------------------
```

And then, deploy the binary to your raspberry pi by,
```shell
$ scp app.pi pi@192.168.31.57:/home/pi
```
`192.168.31.57` is the ip address of my raspberry pi, you need to replace it with yours.

ssh to you raspberry pi, and run the binary.
```shell
$ ssh pi@192.168.31.57
$ ./devices.pi

# or, run it in background
$ nohub ./devices.pi > devices.pi 2>&1 &
```

## App
### [Self-Dirving Car](/app/car)
<img src="img/car.gif" width=80% height=80% />

### [Video Monitor](/app/vmonitor)
<img src="img/vmonitor.gif" width=80% height=80% />

### [Auto-Air](/app/autoair)
<img src="img/auto-air.gif" width=80% height=80% />

### [Auto-Light](/app/autolight)
<img src="img/auto-light.gif" width=80% height=80% />

## [Auto-Fan](/app/autofan)
<img src="img/auto-fan.gif" width=40% height=40% />
