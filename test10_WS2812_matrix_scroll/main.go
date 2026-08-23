// WS2812を35個（横7・縦5、蛇行配線）並べたマトリクスに「テスト」の文字を
// 右から左へ流しながら表示するサンプル。1周ごとに色相を変える。
package main

import (
	"image/color"
	"machine"
	"math"
	"time"

	"tinygo.org/x/drivers/ws2812"
)

const (
	ledPin  = machine.GP15 // WS2812のDINを接続するピン
	numCols = 7            // マトリクスの横方向のLED数
	numRows = 5            // マトリクスの縦方向のLED数
	numLEDs = numCols * numRows

	brightness  = 0.15                   // 明るさ (0.0-1.0)。眩しさ・電流対策で控えめにしている
	scrollDelay = 100 * time.Millisecond // 1列スクロールするごとの待ち時間
	hueStep     = 47.0                   // 1周流れるごとに変化させる色相の量(度)
)

const message = "テスト"

// glyphs は各文字を5列×5行のドットパターンで表現したフォント。
// KanaFive(5×5ドットフォント)を参照して作成した。
// 参考: https://wentwayup.tamaliver.jp/e93615.html
// 各要素が1行分（上から下）を表し、'1'が点灯、'0'が消灯を表す。
var glyphs = map[rune][numRows]string{
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

func main() {
	ledPin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	ws := ws2812.NewWS2812(ledPin)

	buf := buildBuffer(message)

	var leds [numLEDs]color.RGBA
	hue := 0.0

	for {
		on := hsvToRGB(hue, 1.0, brightness)
		off := color.RGBA{}

		for start := -(numCols - 1); start <= len(buf)-1; start++ {
			for c := 0; c < numCols; c++ {
				col := columnAt(buf, start+c)
				for r := 0; r < numRows; r++ {
					idx := ledIndex(r, c)
					if col&(1<<uint(numRows-1-r)) != 0 {
						leds[idx] = on
					} else {
						leds[idx] = off
					}
				}
			}
			ws.WriteColors(leds[:])
			time.Sleep(scrollDelay)
		}

		hue = math.Mod(hue+hueStep, 360)
	}
}

// ledIndex は画面上の行r・列c（ともに0始まり）から、蛇行配線されたLED配列上の
// インデックス(0始まり)を返す。偶数行は左→右、奇数行は右→左に配線されている。
func ledIndex(r, c int) int {
	if r%2 == 0 {
		return r*numCols + c
	}
	return r*numCols + (numCols - 1 - c)
}

// buildBuffer はメッセージ文字列を、1列ごとのビットパターン（下位numRowsビットを使用）
// のスライスに変換する。文字間には1列分の空白を挟む。
func buildBuffer(msg string) []uint8 {
	var cols []uint8
	for i, r := range []rune(msg) {
		if i > 0 {
			cols = append(cols, 0)
		}
		g, ok := glyphs[r]
		if !ok {
			continue
		}
		width := len(g[0])
		for c := 0; c < width; c++ {
			var bits uint8
			for row := 0; row < numRows; row++ {
				if g[row][c] == '1' {
					bits |= 1 << uint(numRows-1-row)
				}
			}
			cols = append(cols, bits)
		}
	}
	return cols
}

// columnAt はバッファの範囲外を空白列(0)として扱うヘルパー。
func columnAt(buf []uint8, idx int) uint8 {
	if idx < 0 || idx >= len(buf) {
		return 0
	}
	return buf[idx]
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
