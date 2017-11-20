package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/splace/joysticks"
	"github.com/tarm/serial"
)

type GamepadOptions struct {
	Port     string `short:"p" long:"port" description:"Port to send and receive data, per example '/dev/ttyUSB0'." default:"/dev/ttyUSB0"`
	Speed    int    `short:"s" long:"speed" description:"Baud rate." default:"115200"`
	MaxValue int    `short:"m" long:"max-value" description:"Max value to send to the receiver" default:"400"`
	Path     string `long:"path" description:"Path where the data file is stored" default:"out.csv"`
}

func (g *GamepadOptions) Execute(args []string) error {
	maxValue := g.MaxValue

	// Initialize serial port
	c := &serial.Config{Name: g.Port, Baud: g.Speed, ReadTimeout: time.Millisecond * 200}
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

				x := mapValues(hpos.X, -1, 1, float32(maxValue), -float32(maxValue))
				y := mapValues(hpos.Y, -1, 1, float32(maxValue), -float32(maxValue))

				rMotor = int(constrain(y+x, -float32(maxValue), float32(maxValue)))
				lMotor = int(constrain(y-x, -float32(maxValue), float32(maxValue)))
			}
		}
	}()

	f, err := os.Create(g.Path)
	if err != nil {
		log.Fatalln("Error creating file", err)
	}

	w := csv.NewWriter(f)
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

		if strings.Contains(string(buf[:n]), "ERROR") {
			log.Println("Communication error")
			continue
		}

		rows := bytes.Split(buf[:n], []byte("|"))
		for _, rawRows := range rows {
			fields := bytes.Split(rawRows, []byte(","))
			rowToWrite := []string{}
			for _, rawField := range fields {
				if len(rawField) != 0 {
					rowToWrite = append(rowToWrite, string(rawField))
				}
			}

			if len(rowToWrite) != 0 {
				if err := w.Write(rowToWrite); err != nil {
					log.Fatal("error writting to csv file", err)
				}
			}
		}

		w.Flush()
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
