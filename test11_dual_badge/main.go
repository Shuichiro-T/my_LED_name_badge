package main

import "time"

// 電子ペーパー(EPD213BWR)と円形ディスプレイ(GC9A01)を同時接続し、
// 電子ペーパーには名前とXアカウント、円形ディスプレイにはicon/myicon.jpgから変換したアイコンを表示する名札プログラム。
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

	for {
		time.Sleep(time.Hour)
	}
}
