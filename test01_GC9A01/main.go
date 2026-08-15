package main

import (
	"image/color"
	"machine"
	"time"

	"tinygo.org/x/drivers/gc9a01"
)

func main() {
	// USBシリアルが安定するまで待つ。tinygo monitor等で接続してから出力を見る。
	time.Sleep(3 * time.Second)
	println("boot: start")

	// [切り分け用に一時変更中] SPI1 + 別ピンに配線を移して同じ症状が出るか確認する。
	// 問題が解決したら元のSPI0構成（コメントアウト部分）に戻すこと。
	machine.SPI1.Configure(machine.SPIConfig{
		Frequency: 1000000, // さらに低速にして信号品質要因を排除
		SCK:       machine.GPIO10,
		SDO:       machine.GPIO11, // MOSI
		SDI:       machine.GPIO8,  // MISO (未使用)
	})
	println("boot: spi configured")

	// CSピンは常時選択状態（Low）に固定する。
	// ドライバは描画時にCSを明示的に操作しないため、起動直後の不定な状態に依存させない。
	csPin := machine.GPIO9
	csPin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	csPin.Low()

	display := gc9a01.New(machine.SPI1, machine.GPIO12, machine.GPIO13, csPin, machine.GPIO14)
	println("boot: calling display.Configure")
	display.Configure(gc9a01.Config{Orientation: gc9a01.HORIZONTAL, Width: 240, Height: 240})
	println("boot: display configured")
	display.EnableBacklight(true)

	red := color.RGBA{255, 0, 0, 255}
	if err := display.FillRectangle(0, 0, 240, 240, red); err != nil {
		println("FillRectangle error:", err.Error())
	}
	println("boot: fill screen done")

	for {
		time.Sleep(time.Hour)
	}
}
