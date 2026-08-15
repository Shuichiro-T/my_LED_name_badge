# test07_EPD213BWR

WeAct Studio 2.13インチ 電子ペーパーモジュール（Model: E0213A179、Drive IC: SSD1680、122×250px、白黒赤(BWR)3色）に「しゅういちろ」と表示する名札プログラムです。TinyGo で実装しています。

- 参考: [WeActStudio.EpaperModule](https://github.com/WeActStudio/WeActStudio.EpaperModule) / `Doc/2.13 Inch Black&Write&Red/英瑞达E0213A179（BWR）.pdf`
- `tinygo.org/x/drivers` には該当パネル用のドライバが無いため、SSD1680を直接SPI制御する自前ドライバを実装している（外部依存なし）。
- 初期化・RAM書き込みシーケンスは、ベンダー提供のMSP430向けサンプル（コメント文字化けあり）ではなく、同じSSD1680・122×250 BWRパネル向けの実績あるArduinoライブラリ [ZinggJM/GxEPD2](https://github.com/ZinggJM/GxEPD2)（`GxEPD2_213_Z98c`）を参考にした。

## 接続（配線）

SPI1 + GPIO10-14。`test02`/`test03_GC9A01` と同じSPI1・ピン範囲に、EPD特有のBUSY入力ピンを加えた構成。

| EPDモジュール側 | Pico 2 側 | 役割 |
|---|---|---|
| VCC | 3V3 (OUT) | 電源 (3.3V) |
| GND | GND | グラウンド |
| SCL / SCK | GPIO10 | SPIクロック |
| SDA / SDI(MOSI) | GPIO11 | SPIデータ (Tx) |
| CS# | GPIO9 | チップセレクト |
| D/C# | GPIO13 | データ/コマンド切替 |
| RES# | GPIO12 | リセット |
| BUSY | GPIO14 | ビジー状態入力（Highの間はコマンド送信禁止） |

配線を変更する場合は [main.go](main.go) 冒頭の `spi`/`sckPin`/`sdoPin`/`csPin`/`resetPin`/`dcPin`/`busyPin` を書き換えること。

## 表示内容・向きについて

- パネルは物理的には縦長（122px幅 × 250px高さ）だが、名札として横書きで読めるよう、ソフトウェア側で90度回転して描画している。
- そのため、**モジュールを縦長のまま実装した場合は90度横に傾けて（横長になるように）見る/装着する**必要がある。逆向きにしたい場合は [main.go](main.go) の `setLogicalPixel` 内の回転変換を書き換えること。
- 白背景に黒文字で「しゅういちろ」、周囲に赤の枠線を描画する（BWR3色対応の確認を兼ねる）。
- 描画は起動時に1回だけフル更新（画面全体書き換え）を行い、その後ディープスリープに入る（電子ペーパーは電源を切っても表示が保持されるため）。

## 必要なツール

- [TinyGo](https://tinygo.org/) (v0.41.1で動作確認)
- Go 1.25以降

## ビルド方法

```sh
cd test07_EPD213BWR
tinygo build -target=pico2 -o /tmp/out.uf2 .   # ビルド確認のみ
```

## 書き込み方法

```sh
tinygo flash -target=pico2 .   # BOOTSELモードで接続した状態で実行
```

## フォントについて

`し`, `ゅ`, `う`, `い`, `ち`, `ろ` の6文字分のみ、40x40pxの1bitビットマップとして [font.go](font.go) に埋め込んでいる（macOSの「ヒラギノ角ゴシック W6」フォントから生成）。他の文字を表示したい場合はグリフデータを追加生成する必要がある。
