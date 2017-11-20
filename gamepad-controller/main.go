package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/splace/joysticks"
	"github.com/tarm/serial"
)

func main() {
	port := flag.String("port", "/dev/ttyUSB0", "Port to send and receive data, per example '/dev/ttyUSB0'.")
	speed := flag.Int("speed", 115200, "Baud rate.")
	maxValue := flag.Int("max-value", 400, "Max value to send to the receiver.")

	// Initialize serial port
	c := &serial.Config{Name: *port, Baud: *speed, ReadTimeout: time.Millisecond * 200}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}

	// Initialize keypad
	device := joysticks.Connect(1)
	if device == nil {
		log.Fatal("no HIDs")
	}
	log.Printf("HID#1:- Buttons:%d, Hats:%d\n", len(device.Buttons), len(device.HatAxes)/2)

	h1move := device.OnMove(1)

	// feed OS events onto the event channels.
	go device.ParcelOutEvents()

	var lMotor int
	var rMotor int

	go func() {
		for {
			select {
			case h := <-h1move:
				hpos := h.(joysticks.CoordsEvent)

				x := mapValues(hpos.X, -1, 1, float32(*maxValue), -float32(*maxValue))
				y := mapValues(hpos.Y, -1, 1, float32(*maxValue), -float32(*maxValue))

				rMotor = int(constrain(y+x, -float32(*maxValue), float32(*maxValue)))
				lMotor = int(constrain(y-x, -float32(*maxValue), float32(*maxValue)))
			}
		}
	}()

	for {
		writeStr := fmt.Sprintf("%d,%d|", lMotor, rMotor)
		log.Println("WRITE", writeStr)
		_, err := s.Write([]byte(writeStr))
		if err != nil {
			log.Println("Error writting to serial port", err)
		}

		buf := make([]byte, 128)
		n, err := s.Read(buf)
		if err != nil && err != io.EOF {
			log.Fatal(err)
		}
		log.Printf("%q", buf[:n])
	}
}

func mapValues(x, inMin, inMax, outMin, outMax float32) float32 {
	return (x-inMin)*(outMax-outMin)/(inMax-inMin) + outMin
}

func constrain(x, a, b float32) float32 {
	if x > a && x < b {
		return x
	}

	if x < a {
		return a
	}

	if x > b {
		return b
	}

	return x
}
