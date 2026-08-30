package main

import (
	"image/color"
	"machine"
	"time"

	"tinygo.org/x/tinyfont"
	"tinygo.org/x/tinyfont/freemono"
)

// WeAct Studio 2.13インチ EPD module (Model: E0213A179, Drive IC: SSD1680, 122x250px, Black/White/Red)
// test07_EPD213BWR の実装をそのまま移植し、表示内容を「名前＋Xアカウント」の2段組に拡張したもの。
// 参考: https://github.com/WeActStudio/WeActStudio.EpaperModule
//       Doc/2.13 Inch Black&Write&Red/英瑞达E0213A179（BWR）.pdf

// パネル物理解像度。RAM上のX方向は8px単位（バイト境界）で、122pxの表示領域に対し128px(16byte)分確保される。
const (
	epdPanelWidthBytes   = 16 // 128px / 8
	epdPanelVisibleWidth = 122
	epdPanelHeight       = 250
	epdLogicalWidth      = epdPanelHeight       // 名札を横向きに使うため90度回転して描画する
	epdLogicalHeight     = epdPanelVisibleWidth // 回転後の縦幅は物理X方向の可視幅に一致させる
)

// SPI1 + GPIO8-14。円形ディスプレイ(SPI0)とは独立したバス・ピンで同時接続する。
var (
	epdSPI      = machine.SPI1
	epdSCKPin   = machine.GPIO10
	epdSDOPin   = machine.GPIO11 // MOSI (SDA)
	epdSDIPin   = machine.GPIO8  // MISO (未使用)
	epdCSPin    = machine.GPIO9
	epdResetPin = machine.GPIO12
	epdDCPin    = machine.GPIO13
	epdBusyPin  = machine.GPIO14 // Busy高進行中はコマンド送信禁止（データシート Note 6-4）
)

var colorEPDRed = color.RGBA{R: 255, A: 255}

// epdBlackBuf: 1=白(発色なし) / 0=黒インク。epdRedBuf: 1=赤インク / 0=赤なし。
// いずれもSSD1680の内蔵RAMと同じレイアウト（1行 epdPanelWidthBytes バイト、epdPanelHeight 行）。
var (
	epdBlackBuf [epdPanelWidthBytes * epdPanelHeight]byte
	epdRedBuf   [epdPanelWidthBytes * epdPanelHeight]byte
)

// epdCanvas は tinyfont から見た描画先。文字は常に赤インクで焼き込む（drivers.Displayer実装）。
type epdCanvas struct{}

func (epdCanvas) Size() (int16, int16) { return int16(epdLogicalWidth), int16(epdLogicalHeight) }

func (epdCanvas) SetPixel(x, y int16, _ color.RGBA) {
	epdSetLogicalPixel(int(x), int(y), false, true)
}

func (epdCanvas) Display() error { return nil }

func epdCmdData(c byte, bs ...byte) {
	epdCSPin.Low()
	epdDCPin.Low()
	epdSPI.Transfer(c)
	if len(bs) > 0 {
		epdDCPin.High()
		for _, b := range bs {
			epdSPI.Transfer(b)
		}
	}
	epdCSPin.High()
}

func epdCmdDataBuf(c byte, buf []byte) {
	epdCSPin.Low()
	epdDCPin.Low()
	epdSPI.Transfer(c)
	epdDCPin.High()
	for _, b := range buf {
		epdSPI.Transfer(b)
	}
	epdCSPin.High()
}

func epdWaitBusy() {
	for epdBusyPin.Get() {
		time.Sleep(10 * time.Millisecond)
	}
}

func epdReset() {
	epdResetPin.High()
	time.Sleep(20 * time.Millisecond)
	epdResetPin.Low()
	time.Sleep(10 * time.Millisecond)
	epdResetPin.High()
	time.Sleep(20 * time.Millisecond)
	epdWaitBusy()
}

// epdSetRamArea はRAMのX/YアドレスウィンドウとRAMアドレスカウンタを設定する。
// ZinggJM/GxEPD2 の GxEPD2_213_Z98c（同じSSD1680・122x250 BWRパネル向けドライバ）の
// _setPartialRamArea を参考に実装（データシートの英語コメントが破損していたベンダーサンプルより信頼できるため）。
func epdSetRamArea(x, y, w, h int) {
	epdCmdData(0x44, byte(x/8), byte((x+w-1)/8))
	epdCmdData(0x45, byte(y%256), byte(y/256), byte((y+h-1)%256), byte((y+h-1)/256))
	epdCmdData(0x4E, byte(x/8))
	epdCmdData(0x4F, byte(y%256), byte(y/256))
}

// epdConfigurePins はSPIバスと各制御ピンを初期化する（ハードウェアリセットは行わない）。
func epdConfigurePins() {
	epdSPI.Configure(machine.SPIConfig{
		Frequency: 4000000,
		SCK:       epdSCKPin,
		SDO:       epdSDOPin,
		SDI:       epdSDIPin,
	})

	epdCSPin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	epdCSPin.High()
	epdResetPin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	epdDCPin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	epdBusyPin.Configure(machine.PinConfig{Mode: machine.PinInput})
}

func epdInitPanel() {
	epdReset()

	epdCmdData(0x12) // SWRESET
	epdWaitBusy()

	epdCmdData(0x01, 0xF9, 0x00, 0x00) // Driver output control: MUX=250-1=0xF9
	epdCmdData(0x11, 0x03)             // Data entry mode: X/Yともにインクリメント
	epdCmdData(0x3C, 0x05)             // Border waveform
	epdCmdData(0x18, 0x80)             // 内蔵温度センサを使用
	epdCmdData(0x21, 0x00, 0x80)       // Display update control（赤RAMの極性設定を含む）

	epdSetRamArea(0, 0, epdPanelVisibleWidth+6, epdPanelHeight) // X方向はバイト境界のため128px分確保
}

func epdSendFramebuffer() {
	epdSetRamArea(0, 0, epdPanelVisibleWidth+6, epdPanelHeight)
	epdCmdDataBuf(0x24, epdBlackBuf[:])
	epdCmdDataBuf(0x26, epdRedBuf[:])
}

func epdUpdateFull() {
	epdCmdData(0x22, 0xF7) // Display Update Sequence Option: OTPからLUTロード＋フル更新
	epdCmdData(0x20)       // Master Activation
	epdWaitBusy()
}

func epdDeepSleep() {
	epdCmdData(0x10, 0x01)
}

// epdSetPixel は物理座標(px, py)へ描画する。px: 0-121 (可視領域), py: 0-249。
func epdSetPixel(px, py int, ink, red bool) {
	if px < 0 || px >= epdPanelVisibleWidth || py < 0 || py >= epdPanelHeight {
		return
	}
	idx := py*epdPanelWidthBytes + px/8
	mask := byte(0x80 >> uint(px%8))
	if ink {
		epdBlackBuf[idx] &^= mask
	} else {
		epdBlackBuf[idx] |= mask
	}
	if red {
		epdRedBuf[idx] |= mask
	} else {
		epdRedBuf[idx] &^= mask
	}
}

// epdSetLogicalPixel は横向き名札用の論理座標(lx: 0-249, ly: 0-121)を物理座標へ90度回転してマッピングする。
func epdSetLogicalPixel(lx, ly int, ink, red bool) {
	px := epdLogicalHeight - 1 - ly
	py := lx
	epdSetPixel(px, py, ink, red)
}

func epdFillWhite() {
	for i := range epdBlackBuf {
		epdBlackBuf[i] = 0xFF
		epdRedBuf[i] = 0x00
	}
}

func epdDrawBorder(thickness int) {
	for ly := 0; ly < epdLogicalHeight; ly++ {
		for lx := 0; lx < epdLogicalWidth; lx++ {
			if lx < thickness || lx >= epdLogicalWidth-thickness || ly < thickness || ly >= epdLogicalHeight-thickness {
				epdSetLogicalPixel(lx, ly, false, true)
			}
		}
	}
}

func epdDrawGlyph(ch rune, lx0, ly0 int) {
	bits := hiragana[ch]
	if bits == nil {
		return
	}
	wb := (glyphSize + 7) / 8
	for row := 0; row < glyphSize; row++ {
		for col := 0; col < glyphSize; col++ {
			b := bits[row*wb+col/8] & (0x80 >> uint(col%8))
			if b != 0 {
				epdSetLogicalPixel(lx0+col, ly0+row, true, false)
			}
		}
	}
}

func epdDrawName(s string, lx0, ly0 int) {
	x := lx0
	for _, ch := range s {
		epdDrawGlyph(ch, x, ly0)
		x += glyphSize
	}
}

// epdDrawBadge は名前(ひらがな、glyphSizeビットマップ)とXアカウント(ASCII、tinyfont)を
// 上下2段に配置して名札の内容をフレームバッファへ描き込む。
func epdDrawBadge(name, handle string) {
	epdFillWhite()
	epdDrawBorder(3)

	// 2段組の合計高さ(名前40px + 間隔12px + Xアカウント約15px)を、盤面の縦幅(122px)に対して中央に配置する。
	const (
		nameGap    = 12
		handleAsc  = 11 // freemono.Regular9pt7b の上端(ベースラインからの高さ)
		handleDesc = 4  // 同フォントの下端(ベースラインからの深さ)
	)
	blockHeight := glyphSize + nameGap + handleAsc + handleDesc
	nameY := (epdLogicalHeight - blockHeight) / 2

	nameWidth := glyphSize * len([]rune(name))
	nameX := (epdLogicalWidth - nameWidth) / 2
	epdDrawName(name, nameX, nameY)

	var canvas epdCanvas
	handleWidth, _ := tinyfont.LineWidth(&freemono.Regular9pt7b, handle)
	handleX := (epdLogicalWidth - int(handleWidth)) / 2
	handleBaselineY := nameY + glyphSize + nameGap + handleAsc
	tinyfont.WriteLine(&canvas, &freemono.Regular9pt7b, int16(handleX), int16(handleBaselineY), handle, colorEPDRed)
}
