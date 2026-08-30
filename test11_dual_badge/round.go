package main

import (
	"image/color"
	"machine"

	"tinygo.org/x/drivers/gc9a01"
)

// GC9A01 240x240円形ディスプレイに、icon/myicon.jpg から変換したビットマップ（icon_bitmap.go）を表示する。
// test02_GC9A01 の接続方法を踏襲しつつ、SPI0側へ移設して電子ペーパー(SPI1)と共存させる。

// SPI0 + GPIO0-7,15。EPD側(SPI1 + GPIO8-14)とはバス・ピンとも独立しているため同時接続できる。
// ピンヘッダーに出ているGPIO0-15,26-29の範囲に収めるため、SPI0のデフォルト(GPIO16/18/19)ではなく
// GPIO4/6/7を使う（RP2350のSPI0はSCK:2/6/18/22, SDO:3/7/19/23, SDI:0/4/16/20のいずれかに限定される）。
var (
	gcSCKPin   = machine.GPIO6
	gcSDOPin   = machine.GPIO7 // MOSI
	gcSDIPin   = machine.GPIO4 // MISO (未使用)
	gcCSPin    = machine.GPIO5
	gcDCPin    = machine.GPIO15
	gcResetPin = machine.GPIO0
	gcBLPin    = machine.GPIO1 // バックライト
)

var (
	colorGCBackground = color.RGBA{R: 255, G: 255, B: 255, A: 255} // 白背景
	colorGCMark       = color.RGBA{R: 255, A: 255}                 // 赤インク
)

// gcIconBuf は xIconBitmap(1bpp) を画面全面分の色配列に展開したもの。
// FillRectangleWithBuffer で一括転送するため、SetPixelを57600回呼ぶより高速に描画できる。
var gcIconBuf [xIconSize * xIconSize]color.RGBA

func gcBuildIconBuffer() {
	const rowBytes = xIconSize / 8
	i := 0
	for y := 0; y < xIconSize; y++ {
		for x := 0; x < xIconSize; x++ {
			b := xIconBitmap[y*rowBytes+x/8]
			if b&(0x80>>uint(x%8)) != 0 {
				gcIconBuf[i] = colorGCMark
			} else {
				gcIconBuf[i] = colorGCBackground
			}
			i++
		}
	}
}

func gcInitAndDrawIcon() {
	machine.SPI0.Configure(machine.SPIConfig{
		Frequency: 20000000,
		SCK:       gcSCKPin,
		SDO:       gcSDOPin,
		SDI:       gcSDIPin,
	})

	// CSピンは常時選択状態（Low）に固定する。
	// ドライバは描画時にCSを明示的に操作しないため、起動直後の不定な状態に依存させない。
	gcCSPin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	gcCSPin.Low()

	display := gc9a01.New(machine.SPI0, gcResetPin, gcDCPin, gcCSPin, gcBLPin)
	display.Configure(gc9a01.Config{Orientation: gc9a01.HORIZONTAL, Width: 240, Height: 240})
	display.EnableBacklight(true)

	gcBuildIconBuffer()
	display.FillRectangleWithBuffer(0, 0, xIconSize, xIconSize, gcIconBuf[:])
	display.Display()
}
