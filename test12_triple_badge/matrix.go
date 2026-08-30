package main

import (
	"image/color"
	"machine"
	"math"
	"time"

	"tinygo.org/x/drivers/ws2812"
)

// WS2812を35個（横7・縦5、蛇行配線）並べたマトリクスに「テスト」の文字を
// 右から左へ流しながら表示する。test10_WS2812_matrix_scroll の実装を移植したもの。
// DINピンはEPD(SPI1)・GC9A01(SPI0)と同時接続するため、test10のGP15から
// 空いているGP2に変更している（配線を変える場合はこの定数を書き換えること）。

const (
	matrixLedPin  = machine.GPIO2 // WS2812のDINを接続するピン
	matrixCols    = 7             // マトリクスの横方向のLED数
	matrixRows    = 5             // マトリクスの縦方向のLED数
	matrixNumLEDs = matrixCols * matrixRows

	matrixBrightness  = 0.08                   // 明るさ (0.0-1.0)。眩しさ・電流対策で控えめにしている
	matrixScrollDelay = 100 * time.Millisecond // 1列スクロールするごとの待ち時間
	matrixHueStep     = 47.0                   // 1周流れるごとに変化させる色相の量(度)
)

const matrixMessage = "テスト"

// matrixGlyphs は各文字を5列×5行のドットパターンで表現したフォント（KanaFive由来、test10と同一）。
var matrixGlyphs = map[rune][matrixRows]string{
	'テ': {
		"01110",
		"00000",
		"11111",
		"00100",
		"00100",
	},
	'ス': {
		"11110",
		"00010",
		"00100",
		"01010",
		"10001",
	},
	'ト': {
		"10000",
		"10000",
		"11100",
		"10011",
		"10000",
	},
}

// matrixRun はマトリクスへのスクロール描画を無限に繰り返す（呼び出し元へは戻らない）。
func matrixRun() {
	matrixLedPin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	ws := ws2812.NewWS2812(matrixLedPin)

	buf := matrixBuildBuffer(matrixMessage)

	var leds [matrixNumLEDs]color.RGBA
	hue := 0.0

	for {
		on := matrixHSVToRGB(hue, 1.0, matrixBrightness)
		off := color.RGBA{}

		for start := -(matrixCols - 1); start <= len(buf)-1; start++ {
			for c := 0; c < matrixCols; c++ {
				col := matrixColumnAt(buf, start+c)
				for r := 0; r < matrixRows; r++ {
					idx := matrixLedIndex(r, c)
					if col&(1<<uint(matrixRows-1-r)) != 0 {
						leds[idx] = on
					} else {
						leds[idx] = off
					}
				}
			}
			ws.WriteColors(leds[:])
			time.Sleep(matrixScrollDelay)
		}

		hue = math.Mod(hue+matrixHueStep, 360)
	}
}

// matrixLedIndex は画面上の行r・列c（ともに0始まり）から、蛇行配線されたLED配列上の
// インデックス(0始まり)を返す。偶数行は左→右、奇数行は右→左に配線されている。
func matrixLedIndex(r, c int) int {
	if r%2 == 0 {
		return r*matrixCols + c
	}
	return r*matrixCols + (matrixCols - 1 - c)
}

// matrixBuildBuffer はメッセージ文字列を、1列ごとのビットパターン（下位matrixRowsビットを使用）
// のスライスに変換する。文字間には1列分の空白を挟む。
func matrixBuildBuffer(msg string) []uint8 {
	var cols []uint8
	for i, r := range []rune(msg) {
		if i > 0 {
			cols = append(cols, 0)
		}
		g, ok := matrixGlyphs[r]
		if !ok {
			continue
		}
		width := len(g[0])
		for c := 0; c < width; c++ {
			var bits uint8
			for row := 0; row < matrixRows; row++ {
				if g[row][c] == '1' {
					bits |= 1 << uint(matrixRows-1-row)
				}
			}
			cols = append(cols, bits)
		}
	}
	return cols
}

// matrixColumnAt はバッファの範囲外を空白列(0)として扱うヘルパー。
func matrixColumnAt(buf []uint8, idx int) uint8 {
	if idx < 0 || idx >= len(buf) {
		return 0
	}
	return buf[idx]
}

// matrixHSVToRGB は HSV (h: 0-360, s,v: 0.0-1.0) を WS2812 用の RGBA に変換する。
func matrixHSVToRGB(h, s, v float64) color.RGBA {
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
