![](img/go-devices.png)

# rpi-devices 
[![Build Status](https://travis-ci.org/shanghuiyang/rpi-devices.svg?branch=master)](https://travis-ci.org/shanghuiyang/rpi-devices)

rpi-devices let you drive the devices using a raspberry pi in golang, and push your data onto an IoT cloud platform for visualizing.

The following devices had been implemented in the current version, and a device interface was designed to let you add new devices easily.


  ![](img/led.jpg)   ![](img/temp.jpg)   ![](img/relay.jpg)   ![](img/step-motor.jpg)   ![](img/gps.jpg)

         LED        Temperature        Relay        Step-Motor         GPS


  ![](img/buzzer.jpg)   ![](img/voice.jpg)   ![](img/hc-sr04.jpg)   ![](img/l298n.jpg)   ![](img/sg90.jpg)

	  Buzzer          Voice         Ultrasonic        Motor        Servo-Motor
	                  Sensor          Sensor          Dirver

  ![](img/dht11.jpg)   ![](img/infared.jpg)   ![](img/oled.jpg)   ![](img/pms7003.jpg)

	Temp&Humidity    Infared           OLED         Air Quality
   	Sensor       Sensor           Display         (PM2.5)

## Visualize Data
The data from devices can be pushed to an IoT cloud platform for visualizing. rpi-devices designed an interface of IoT cloud, and implemented the interface for [WSN](http://www.wsncloud.com/) cloud and [OneNET](https://open.iot.10086.cn/) cloud. Both of them are free for developers. You can implement the interface for new IoT cloud and add it to the framework easily.

* [WSN](http://www.wsncloud.com/)
    
    visualize the temperature
	
	<img src="img/temp-vis.png" width=70% height=70% />
* [OneNET](https://open.iot.10086.cn/)

    visualize the gps locations

	<img src="img/gps.gif" width=30% height=30% />

## Usage

It is very easy to cross-compile and deploy for golang. It is an example that compiles the binary for raspberry pi on MacOS.
```shell
$ CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=6 go build -o devices.pi main.go
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
$ scp devices.pi pi@192.168.31.57:/home/pi
```
`192.168.31.57` is the ip address of my raspberry pi, you need to replace it with yours.

ssh to you raspberry pi, and run the binary.
```shell
$ ssh pi@192.168.31.57
$ ./devices.pi

# or, run it in background
$ nohub ./devices.pi > devices.pi 2>&1 &
```

## Examples

### LED
```go
package main

import (
	"fmt"
	"log"

	"github.com/shanghuiyang/rpi-devices/dev"
	"github.com/stianeikeland/go-rpio"
)

const (
	p12 = 16 // led
)

func main() {
	if err := rpio.Open(); err != nil {
		log.Fatalf("failed to open rpio, error: %v", err)
		return
	}
	defer rpio.Close()

	led := dev.NewLed(p12)

	var op string
	for {
		fmt.Printf(">>op: ")
		if n, err := fmt.Scanf("%s", &op); n != 1 || err != nil {
			log.Printf("invalid operator, error: %v", err)
			continue
		}
		switch op {
		case "on":
			led.On()
		case "off":
			led.Off()
		case "blink":
			led.Blink(5, 100)
		case "fade":
			led.Fade(3)
		case "q":
			log.Printf("quit\n")
			return
		default:
			fmt.Printf("invalid operator, should be: on, off, blink or q\n")
		}
	}
}
```

### Temperature
```go
package main

import (
	"fmt"

	"github.com/shanghuiyang/rpi-devices/dev"
)

func main() {
	t := dev.NewTemperature()
	c, err := t.GetTemperature()
	if err != nil {
		fmt.Printf("failed to get temperature, error: %v", err)
		return
	}
	fmt.Printf("current temperature: %v", c)
}
```

### Relay
```go
package main

import (
	"fmt"
	"log"

	"github.com/shanghuiyang/rpi-devices/dev"
	"github.com/stianeikeland/go-rpio"
)

const (
	p7 = 7 // relay
)

func main() {
	if err := rpio.Open(); err != nil {
		log.Fatalf("failed to open rpio, error: %v", err)
		return
	}
	defer rpio.Close()

	r := dev.NewRelay(p7)
	var op string
	for {
		fmt.Printf(">>op: ")
		if n, err := fmt.Scanf("%s", &op); n != 1 || err != nil {
			log.Printf("invalid operator, error: %v", err)
			continue
		}
		switch op {
		case "on":
			r.On()
		case "off":
			r.Off()
		case "q":
			log.Printf("done\n")
			return
		default:
			fmt.Printf("invalid operator, should be: on, off or q\n")
		}
	}
}
```

### Step-Motor
```go
package main

import (
	"fmt"
	"log"

	"github.com/shanghuiyang/rpi-devices/dev"
	"github.com/stianeikeland/go-rpio"
)

const (
	p8  = 8  // in1 for step motor
	p25 = 25 // in2 for step motor
	p24 = 24 // in3 for step motor
	p23 = 23 // in4 for step motor
)

func main() {
	if err := rpio.Open(); err != nil {
		log.Fatalf("failed to open rpio, error: %v", err)
		return
	}
	defer rpio.Close()

	m := dev.NewStepMotor(p8, p25, p24, p23)
	log.Printf("step motor is ready for service\n")

	var angle float32
	for {
		fmt.Printf(">>op: ")
		if n, err := fmt.Scanf("%f", &angle); n != 1 || err != nil {
			log.Printf("invalid angle, error: %v", err)
			continue
		}
		if angle == 0 {
			break
		}
		m.Roll(angle)
	}
	log.Printf("step motor stop service\n")
}
```

### GPS
```go
package main

import (
	"log"

	"github.com/shanghuiyang/rpi-devices/dev"
)

func main() {
	g := dev.NewGPS()
	pt, err := g.Loc()
	if err != nil {
		log.Printf("failed, error: %v", err)
		return
	}
	log.Printf("%v", pt)
	g.Close()
}
```

### Distance
```go
package main

import (
	"fmt"

	"github.com/shanghuiyang/rpi-devices/dev"
	"github.com/stianeikeland/go-rpio"
)

const (
	pinTrig = 21
	pinEcho = 26
)

func main() {

	if err := rpio.Open(); err != nil {
		return
	}
	h := dev.NewHCSR04(pinTrig, pinEcho)
	d := h.Dist()
	fmt.Printf("d: %v\n", d)

	rpio.Close()
}
```

## App
### [Self-Dirving Car](/app/car)
<img src="img/car.gif" width=80% height=80% />

### [Auto-Light](/app/autolight)
<img src="img/auto-light.gif" width=80% height=80% />

## [Auto-Fan](/app/autofan)
<img src="img/auto-fan.gif" width=40% height=40% />
