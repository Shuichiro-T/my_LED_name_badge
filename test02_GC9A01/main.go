package main

import (
	"image/color"
	"machine"
	"math"
	"time"

	"tinygo.org/x/drivers/gc9a01"
	"tinygo.org/x/tinydraw"
	"tinygo.org/x/tinyfont"
	"tinygo.org/x/tinyfont/freemono"
)

// 起動時に埋め込む初期時刻。書き込み時刻に合わせて書き換えること。
var initialTime = time.Date(2026, 6, 20, 12, 0, 0, 0, time.UTC)

const (
	centerX     = 120
	centerY     = 120
	faceRadius  = 118
	hourHandLen = 60
	minHandLen  = 95
	secHandLen  = 105
)

var (
	colorBG    = color.RGBA{0, 0, 0, 255}
	colorFace  = color.RGBA{60, 60, 60, 255}
	colorTick  = color.RGBA{200, 200, 200, 255}
	colorHour  = color.RGBA{255, 255, 255, 255}
	colorMin   = color.RGBA{255, 255, 255, 255}
	colorSec   = color.RGBA{255, 60, 60, 255}
	colorLabel = color.RGBA{255, 255, 255, 255}
)

func main() {
	bootTime := time.Now()

	// test01_GC9A01 に合わせたピン配置。
	machine.SPI1.Configure(machine.SPIConfig{
		Frequency: 20000000, // ブレッドボード配線での通信不安定を避けるため低めに設定
		SCK:       machine.GPIO10,
		SDO:       machine.GPIO11, // MOSI
		SDI:       machine.GPIO8,  // MISO (未使用)
	})

	// CSピンは常時選択状態（Low）に固定する。
	// ドライバは描画時にCSを明示的に操作しないため、起動直後の不定な状態に依存させない。
	csPin := machine.GPIO9
	csPin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	csPin.Low()

	display := gc9a01.New(machine.SPI1, machine.GPIO12, machine.GPIO13, csPin, machine.GPIO14)
	display.Configure(gc9a01.Config{Orientation: gc9a01.HORIZONTAL, Width: 240, Height: 240})
	display.EnableBacklight(true)

	for {
		now := initialTime.Add(time.Since(bootTime))
		drawClock(&display, now)
		time.Sleep(time.Second)
	}
}

func drawClock(display *gc9a01.Device, now time.Time) {
	display.FillScreen(colorBG)
	tinydraw.Circle(display, centerX, centerY, faceRadius, colorFace)
	drawTicks(display)
	drawLabels(display)
	drawHands(display, now)
	display.Display()
}

func drawTicks(display *gc9a01.Device) {
	for i := 0; i < 60; i++ {
		angle := float64(i) * 6 * math.Pi / 180
		outer := float64(faceRadius)
		var inner float64
		if i%5 == 0 {
			inner = outer - 14
		} else {
			inner = outer - 6
		}
		x0 := centerX + int16(inner*math.Sin(angle))
		y0 := centerY - int16(inner*math.Cos(angle))
		x1 := centerX + int16(outer*math.Sin(angle))
		y1 := centerY - int16(outer*math.Cos(angle))
		tinydraw.Line(display, x0, y0, x1, y1, colorTick)
	}
}

func drawLabels(display *gc9a01.Device) {
	type label struct {
		text string
		hour int
	}
	labels := []label{{"12", 0}, {"3", 3}, {"6", 6}, {"9", 9}}
	radius := float64(faceRadius - 28)
	for _, l := range labels {
		angle := float64(l.hour) * 30 * math.Pi / 180
		x := centerX + int16(radius*math.Sin(angle))
		y := centerY - int16(radius*math.Cos(angle))
		w, _ := tinyfont.LineWidth(&freemono.Bold12pt7b, l.text)
		tinyfont.WriteLine(display, &freemono.Bold12pt7b, x-int16(w)/2, y+6, l.text, colorLabel)
	}
}

func drawHands(display *gc9a01.Device, now time.Time) {
	hour := float64(now.Hour() % 12)
	minute := float64(now.Minute())
	second := float64(now.Second())

	hourAngle := (hour+minute/60)*30 * math.Pi / 180
	minAngle := (minute+second/60)*6 * math.Pi / 180
	secAngle := second * 6 * math.Pi / 180

	drawHand(display, hourAngle, hourHandLen, colorHour)
	drawHand(display, minAngle, minHandLen, colorMin)
	drawHand(display, secAngle, secHandLen, colorSec)
	tinydraw.FilledCircle(display, centerX, centerY, 4, colorSec)
}

func drawHand(display *gc9a01.Device, angle float64, length int16, c color.RGBA) {
	x1 := centerX + int16(float64(length)*math.Sin(angle))
	y1 := centerY - int16(float64(length)*math.Cos(angle))
	tinydraw.Line(display, centerX, centerY, x1, y1, c)
}
