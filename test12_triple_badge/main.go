package main

import "time"

// 電子ペーパー(EPD213BWR)・円形ディスプレイ(GC9A01)・WS2812 LEDマトリクスを
// 1台のPico 2に同時接続し、それぞれ役割分担して表示するプロトタイプデモ。
// test11_dual_badge（電子ペーパー＋円形ディスプレイ）をベースに、
// test10_WS2812_matrix_scroll のLEDマトリクス表示を追加したもの。
const (
	badgeName   = "しゅういちろ"
	badgeHandle = "@shucho0103"
)

func main() {
	time.Sleep(2 * time.Second) // 電源投入直後の安定待ち

	epdConfigurePins()
	epdDrawBadge(badgeName, badgeHandle)
	epdInitPanel()
	epdSendFramebuffer()
	epdUpdateFull()
	epdDeepSleep() // 電子ペーパーは電源を切っても表示が保持されるためスリープに入れる

	gcInitAndDrawIcon()

	matrixRun() // LEDマトリクスのスクロール表示を無限に繰り返す（ここから戻らない）
}
