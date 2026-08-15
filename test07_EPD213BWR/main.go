package main

import (
	"machine"
	"time"
)

// WeAct Studio 2.13インチ EPD module (Model: E0213A179, Drive IC: SSD1680, 122x250px, Black/White/Red)
// 参考: https://github.com/WeActStudio/WeActStudio.EpaperModule
//       Doc/2.13 Inch Black&Write&Red/英瑞达E0213A179（BWR）.pdf
//       tinygo.org/x/drivers に該当ドライバが無いため、SSD1680を直接SPI制御する。

// パネル物理解像度。RAM上のX方向は8px単位（バイト境界）で、122pxの表示領域に対し128px(16byte)分確保される。
const (
	panelWidthBytes   = 16 // 128px / 8
	panelVisibleWidth = 122
	panelHeight       = 250
	logicalWidth      = panelHeight       // 名札を横向きに使うため90度回転して描画する
	logicalHeight     = panelVisibleWidth // 回転後の縦幅は物理X方向の可視幅に一致させる
)

// SPI1 + GPIO10-14。test02/03_GC9A01 と同じSPI1・ピン範囲に、EPD特有のBUSY入力を追加した構成。
var (
	spi      = machine.SPI1
	sckPin   = machine.GPIO10
	sdoPin   = machine.GPIO11 // MOSI (SDA)
	sdiPin   = machine.GPIO8  // MISO (未使用)
	csPin    = machine.GPIO9
	resetPin = machine.GPIO12
	dcPin    = machine.GPIO13
	busyPin  = machine.GPIO14 // Busy高進行中はコマンド送信禁止（データシート Note 6-4）
)

// blackBuf: 1=白(発色なし) / 0=黒インク。redBuf: 1=赤インク / 0=赤なし。
// いずれもSSD1680の内蔵RAMと同じレイアウト（1行 panelWidthBytes バイト、panelHeight 行）。
var (
	blackBuf [panelWidthBytes * panelHeight]byte
	redBuf   [panelWidthBytes * panelHeight]byte
)

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

func cmdDataBuf(c byte, buf []byte) {
	csPin.Low()
	dcPin.Low()
	spi.Transfer(c)
	dcPin.High()
	for _, b := range buf {
		spi.Transfer(b)
	}
	csPin.High()
}

func waitBusy() {
	for busyPin.Get() {
		time.Sleep(10 * time.Millisecond)
	}
}

func reset() {
	resetPin.High()
	time.Sleep(20 * time.Millisecond)
	resetPin.Low()
	time.Sleep(10 * time.Millisecond)
	resetPin.High()
	time.Sleep(20 * time.Millisecond)
	waitBusy()
}

// setRamArea はRAMのX/YアドレスウィンドウとRAMアドレスカウンタを設定する。
// ZinggJM/GxEPD2 の GxEPD2_213_Z98c（同じSSD1680・122x250 BWRパネル向けドライバ）の
// _setPartialRamArea を参考に実装（データシートの英語コメントが破損していたベンダーサンプルより信頼できるため）。
func setRamArea(x, y, w, h int) {
	cmdData(0x44, byte(x/8), byte((x+w-1)/8))
	cmdData(0x45, byte(y%256), byte(y/256), byte((y+h-1)%256), byte((y+h-1)/256))
	cmdData(0x4E, byte(x/8))
	cmdData(0x4F, byte(y%256), byte(y/256))
}

func initDisplay() {
	reset()

	cmdData(0x12) // SWRESET
	waitBusy()

	cmdData(0x01, 0xF9, 0x00, 0x00) // Driver output control: MUX=250-1=0xF9
	cmdData(0x11, 0x03)             // Data entry mode: X/Yともにインクリメント
	cmdData(0x3C, 0x05)             // Border waveform
	cmdData(0x18, 0x80)             // 内蔵温度センサを使用
	cmdData(0x21, 0x00, 0x80)       // Display update control（赤RAMの極性設定を含む）

	setRamArea(0, 0, panelVisibleWidth+6, panelHeight) // X方向はバイト境界のため128px分確保
}

func sendFramebuffer() {
	setRamArea(0, 0, panelVisibleWidth+6, panelHeight)
	cmdDataBuf(0x24, blackBuf[:])
	cmdDataBuf(0x26, redBuf[:])
}

func updateFull() {
	cmdData(0x22, 0xF7) // Display Update Sequence Option: OTPからLUTロード＋フル更新
	cmdData(0x20)       // Master Activation
	waitBusy()
}

func deepSleep() {
	cmdData(0x10, 0x01)
}

// setPixel は物理座標(px, py)へ描画する。px: 0-121 (可視領域), py: 0-249。
func setPixel(px, py int, ink, red bool) {
	if px < 0 || px >= panelVisibleWidth || py < 0 || py >= panelHeight {
		return
	}
	idx := py*panelWidthBytes + px/8
	mask := byte(0x80 >> uint(px%8))
	if ink {
		blackBuf[idx] &^= mask
	} else {
		blackBuf[idx] |= mask
	}
	if red {
		redBuf[idx] |= mask
	} else {
		redBuf[idx] &^= mask
	}
}

// setLogicalPixel は横向き名札用の論理座標(lx: 0-249, ly: 0-121)を物理座標へ90度回転してマッピングする。
func setLogicalPixel(lx, ly int, ink, red bool) {
	px := logicalHeight - 1 - ly
	py := lx
	setPixel(px, py, ink, red)
}

func fillWhite() {
	for i := range blackBuf {
		blackBuf[i] = 0xFF
		redBuf[i] = 0x00
	}
}

func drawBorder(thickness int) {
	for ly := 0; ly < logicalHeight; ly++ {
		for lx := 0; lx < logicalWidth; lx++ {
			if lx < thickness || lx >= logicalWidth-thickness || ly < thickness || ly >= logicalHeight-thickness {
				setLogicalPixel(lx, ly, false, true)
			}
		}
	}
}

func drawGlyph(ch rune, lx0, ly0 int) {
	bits := hiragana[ch]
	if bits == nil {
		return
	}
	wb := (glyphSize + 7) / 8
	for row := 0; row < glyphSize; row++ {
		for col := 0; col < glyphSize; col++ {
			b := bits[row*wb+col/8] & (0x80 >> uint(col%8))
			if b != 0 {
				setLogicalPixel(lx0+col, ly0+row, true, false)
			}
		}
	}
}

func drawText(s string, lx0, ly0 int) {
	x := lx0
	for _, ch := range s {
		drawGlyph(ch, x, ly0)
		x += glyphSize
	}
}

func main() {
	time.Sleep(2 * time.Second) // 電源投入直後の安定待ち

	spi.Configure(machine.SPIConfig{
		Frequency: 4000000,
		SCK:       sckPin,
		SDO:       sdoPin,
		SDI:       sdiPin,
	})

	csPin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	csPin.High()
	resetPin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	dcPin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	busyPin.Configure(machine.PinConfig{Mode: machine.PinInput})

	const text = "しゅういちろ"
	textWidth := glyphSize * len([]rune(text))
	startX := (logicalWidth - textWidth) / 2
	startY := (logicalHeight - glyphSize) / 2

	fillWhite()
	drawBorder(3)
	drawText(text, startX, startY)

	initDisplay()
	sendFramebuffer()
	updateFull()
	deepSleep()

	for {
		time.Sleep(time.Hour)
	}
}
