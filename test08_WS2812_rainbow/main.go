// WS2812を9個直列に接続し、虹色のグラデーションを流れるように点灯させるサンプル。
package main

import (
	"image/color"
	"machine"
	"math"
	"time"

	"tinygo.org/x/drivers/ws2812"
)

const (
	ledPin     = machine.GP15 // WS2812のDINを接続するピン
	numLEDs    = 9
	brightness = 0.2 // 明るさ (0.0-1.0)。電流・眩しさ対策で抑えめにしている
)

func main() {
	ledPin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	ws := ws2812.NewWS2812(ledPin)

	var leds [numLEDs]color.RGBA
	var hueOffset float64

	for {
		for i := 0; i < numLEDs; i++ {
			hue := math.Mod(hueOffset+float64(i)*(360.0/numLEDs), 360)
			leds[i] = hsvToRGB(hue, 1.0, brightness)
		}
		ws.WriteColors(leds[:])

		hueOffset = math.Mod(hueOffset+2.0, 360)
		time.Sleep(30 * time.Millisecond)
	}
}

// hsvToRGB は HSV (h: 0-360, s,v: 0.0-1.0) を WS2812 用の RGBA に変換する。
func hsvToRGB(h, s, v float64) color.RGBA {
	c := v * s
	x := c * (1 - math.Abs(math.Mod(h/60, 2)-1))
	m := v - c

	var r, g, b float64
	switch {
	case h < 60:
		r, g, b = c, x, 0
	case h < 120:
		r, g, b = x, c, 0
	case h < 180:
		r, g, b = 0, c, x
	case h < 240:
		r, g, b = 0, x, c
	case h < 300:
		r, g, b = x, 0, c
	default:
		r, g, b = c, 0, x
	}

	return color.RGBA{
		R: uint8((r + m) * 255),
		G: uint8((g + m) * 255),
		B: uint8((b + m) * 255),
		A: 0xff,
	}
}
