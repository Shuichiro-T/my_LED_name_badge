// WS2812を35個直列に接続し、赤い光が先頭から末尾へ移動するサンプル。
package main

import (
	"image/color"
	"machine"
	"time"

	"tinygo.org/x/drivers/ws2812"
)

const (
	ledPin  = machine.GP15 // WS2812のDINを接続するピン
	numLEDs = 35
)

var brightness = 0.1 // 明るさ (0.0-1.0)。電流・眩しさ対策で控えめにしている

func main() {
	ledPin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	ws := ws2812.NewWS2812(ledPin)

	off := color.RGBA{}
	red := color.RGBA{R: uint8(255.0 * brightness), A: 0xff}

	var leds [numLEDs]color.RGBA

	for {
		for pos := 0; pos < numLEDs; pos++ {
			for i := range leds {
				leds[i] = off
			}
			leds[pos] = red
			ws.WriteColors(leds[:])
			time.Sleep(50 * time.Millisecond)
		}
	}
}
