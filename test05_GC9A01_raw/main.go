package main

import (
	"machine"
	"time"
)

// GC9A01 レジスタ定義（tinygo.org/x/drivers/gc9a01 と同じ値）。
const (
	regSLPOUT   = 0x11
	regINVON    = 0x21
	regDISPON   = 0x29
	regCASET    = 0x2A
	regRASET    = 0x2B
	regRAMWR    = 0x2C
	regTEON     = 0x35
	regMADCTR   = 0x36
	regCOLMOD   = 0x3A
	regDISFNCTL = 0xB6
	regFRMCTL   = 0xE8
	regPWCTR3   = 0xC3
	regPWCTR4   = 0xC4
	regINTEN1   = 0xFE
	regGMSET1   = 0xF0
	regGMSET2   = 0xF1
	regGMSET3   = 0xF2
	regGMSET4   = 0xF3
)

// test03で導通確認済みの配線（SPI1 + GPIO10-14）。BLK端子は無い（常時点灯）モジュールなので未使用。
var (
	spi      = machine.SPI1
	sckPin   = machine.GPIO10
	sdoPin   = machine.GPIO11 // MOSI
	sdiPin   = machine.GPIO8  // MISO (未使用)
	csPin    = machine.GPIO9
	resetPin = machine.GPIO12
	dcPin    = machine.GPIO13
)

var bootTime time.Time

func logStep(msg string) {
	println(msg, int64(time.Since(bootTime)/time.Millisecond), "ms")
}

// cmd はコマンドバイトのみをCSフレーム内で送信する（後続のdata()呼び出しが無い単発コマンド用）。
func cmd(b byte) {
	csPin.Low()
	dcPin.Low()
	spi.Transfer(b)
	csPin.High()
}

// cmdData はコマンド＋パラメータを1つのCSフレーム（CS Low〜High）として送信する。
func cmdData(c byte, bs ...byte) {
	csPin.Low()
	dcPin.Low()
	spi.Transfer(c)
	if len(bs) > 0 {
		dcPin.High()
		for _, b := range bs {
			spi.Transfer(b)
		}
	}
	csPin.High()
}

func reset() {
	logStep("reset: RES High")
	resetPin.High()
	time.Sleep(100 * time.Millisecond)
	logStep("reset: RES Low")
	resetPin.Low()
	time.Sleep(100 * time.Millisecond)
	logStep("reset: RES High")
	resetPin.High()
	time.Sleep(100 * time.Millisecond)
	logStep("reset: done")
}

func initDisplay() {
	reset()

	logStep("init: power block start")
	cmdData(0xEF)
	cmdData(0xEB, 0x14)

	cmdData(regINTEN1)
	cmdData(0xEF)

	cmdData(0xEB, 0x14)

	cmdData(0x84, 0x40)
	cmdData(0x85, 0xFF)
	cmdData(0x86, 0xFF)
	cmdData(0x87, 0xFF)
	cmdData(0x88, 0x0A)
	cmdData(0x89, 0x21)
	cmdData(0x8A, 0x00)
	cmdData(0x8B, 0x80)
	cmdData(0x8C, 0x01)
	cmdData(0x8D, 0x01)
	cmdData(0x8E, 0xFF)
	cmdData(0x8F, 0xFF)

	logStep("init: power block done")

	cmdData(regDISFNCTL, 0x00, 0x20)
	cmdData(regMADCTR, 0x08)
	cmdData(regCOLMOD, 0x05)

	logStep("init: MADCTR/COLMOD done")

	cmdData(0x90, 0x08, 0x08, 0x08, 0x08)
	cmdData(0xBD, 0x06)
	cmdData(0xBC, 0x00)
	cmdData(0xFF, 0x60, 0x01, 0x04)
	cmdData(regPWCTR3, 0x13)
	cmdData(regPWCTR4, 0x13)
	cmdData(0xC9, 0x22)
	cmdData(0xBE, 0x11)
	cmdData(0xE1, 0x10, 0x0E)
	cmdData(0xDF, 0x21, 0x0c, 0x02)

	cmdData(regGMSET1, 0x45, 0x09, 0x08, 0x08, 0x26, 0x2A)
	cmdData(regGMSET2, 0x43, 0x70, 0x72, 0x36, 0x37, 0x6F)
	cmdData(regGMSET3, 0x45, 0x09, 0x08, 0x08, 0x26, 0x2A)
	cmdData(regGMSET4, 0x43, 0x70, 0x72, 0x36, 0x37, 0x6F)

	logStep("init: gamma set done")

	cmdData(0xED, 0x1B, 0x0B)
	cmdData(0xAE, 0x77)
	cmdData(0xCD, 0x63)
	cmdData(0x70, 0x07, 0x07, 0x04, 0x0E, 0x0F, 0x09, 0x07, 0x08, 0x03)
	cmdData(regFRMCTL, 0x34)
	cmdData(0x62, 0x18, 0x0D, 0x71, 0xED, 0x70, 0x70, 0x18, 0x0F, 0x71, 0xEF, 0x70, 0x70)
	cmdData(0x63, 0x18, 0x11, 0x71, 0xF1, 0x70, 0x70, 0x18, 0x13, 0x71, 0xF3, 0x70, 0x70)
	cmdData(0x64, 0x28, 0x29, 0xF1, 0x01, 0xF1, 0x00, 0x07)
	cmdData(0x66, 0x3C, 0x00, 0xCD, 0x67, 0x45, 0x45, 0x10, 0x00, 0x00, 0x00)
	cmdData(0x67, 0x00, 0x3C, 0x00, 0x00, 0x00, 0x01, 0x54, 0x10, 0x32, 0x98)
	cmdData(0x74, 0x10, 0x85, 0x80, 0x00, 0x00, 0x4E, 0x00)
	cmdData(0x98, 0x3e, 0x07)

	cmdData(regTEON)
	cmdData(0x21)

	logStep("init: before SLPOUT")
	cmdData(regSLPOUT)
	time.Sleep(120 * time.Millisecond)
	logStep("init: before DISPON")
	cmdData(regDISPON)
	time.Sleep(20 * time.Millisecond)
	logStep("init: after DISPON")
}

func fillRed() {
	const w, h = 240, 240

	logStep("fill: CASET/RASET start")
	cmdData(regCASET, 0x00, 0x00, byte((w-1)>>8), byte(w-1))
	cmdData(regRASET, 0x00, 0x00, byte((h-1)>>8), byte(h-1))

	logStep("fill: RAMWR start")
	csPin.Low()
	dcPin.Low()
	spi.Transfer(regRAMWR)
	dcPin.High()
	// RGB565 で赤一色 = 0xF800。RAMWR全体を1つのCSフレームとして送信する。
	quarter := (w * h) / 4
	for i := 0; i < w*h; i++ {
		spi.Transfer(0xF8)
		spi.Transfer(0x00)
		if i > 0 && i%quarter == 0 {
			logStep("fill: progress")
		}
	}
	csPin.High()
	logStep("fill: RAMWR done")
}

func main() {
	bootTime = time.Now()
	time.Sleep(3 * time.Second)
	logStep("boot: start")

	spi.Configure(machine.SPIConfig{
		Frequency: 1000000,
		SCK:       sckPin,
		SDO:       sdoPin,
		SDI:       sdiPin,
	})
	logStep("boot: spi configured")

	csPin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	csPin.High() // アイドル時はHigh。各コマンド送信時のみcmd()/cmdData()内でLowにする。
	resetPin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	dcPin.Configure(machine.PinConfig{Mode: machine.PinOutput})

	initDisplay()
	logStep("boot: display initialized")

	fillRed()
	logStep("boot: fill red done")

	for {
		time.Sleep(time.Hour)
	}
}
