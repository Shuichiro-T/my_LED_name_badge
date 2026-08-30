# test12_triple_badge

電子ペーパー（test07_EPD213BWR）・円形ディスプレイ（test02_GC9A01系）・WS2812 LEDマトリクス（test10_WS2812_matrix_scroll）を
1台のPico 2に同時接続し、それぞれ同時に動かすプロトタイプデモです。TinyGo で実装しています。

`test11_dual_badge`（電子ペーパー＋円形ディスプレイ）をベースに、`test10_WS2812_matrix_scroll` のLEDマトリクス表示を追加したものです。

- 電子ペーパー（WeAct Studio 2.13インチ EPD213BWR, SSD1680, 122×250px, 白黒赤）: 名前「しゅういちろ」とXアカウント「@shucho0103」を表示。
- 円形ディスプレイ（GC9A01, 240×240px）: [icon/myicon.jpg](icon/myicon.jpg) を変換したアイコン画像を表示。
- WS2812マトリクス（35個、横7×縦5、蛇行配線）: 「テスト」の文字を右から左へ流しながら表示。1周ごとに色相を変える。

電子ペーパーと円形ディスプレイは起動時に1回描画すれば表示が保持される（電子ペーパーは焼き付け式、円形ディスプレイはGRAMが保持）ため、
起動処理の最後にLEDマトリクスのスクロールループ（`matrixRun`、無限ループで戻らない）を実行することで3つを同時に表示させている。

## 接続（配線）

3つのデバイスは独立したSPIバス・ピン（LEDマトリクスはSPIを使わずGPIO1本）を使うため、同時に接続しても干渉しない。
ピンヘッダーに出ているGPIO0-15,26-29の範囲に収まるよう選んでいる。

### 電子ペーパー側（SPI1 + GPIO8-14）test07_EPD213BWR・test11_dual_badgeと同一の配線

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

### 円形ディスプレイ側（SPI0 + GPIO0,1,4-7,15）test11_dual_badgeと同一の配線

| GC9A01モジュール側 | Pico 2 側 | 役割 |
|---|---|---|
| VCC | 3V3 (OUT) | 電源 (3.3V) |
| GND | GND | グラウンド |
| SCL / SCK | GPIO6 | SPIクロック |
| SDA / SDI(MOSI) | GPIO7 | SPIデータ (Tx) |
| CS | GPIO5 | チップセレクト |
| DC | GPIO15 | データ/コマンド切替 |
| RST | GPIO0 | リセット |
| BL | GPIO1 | バックライト |

### WS2812マトリクス側（GPIO2）

test10_WS2812_matrix_scroll ではDINをGP15に接続していたが、GP15は円形ディスプレイのDCピンと競合するため、
GPIO0-15の範囲で空いている **GPIO2** に変更している。

| WS2812テープ側 | Pico 2 側 | 役割 |
|---|---|---|
| VCC (5V or 3V3) | 5V または 3V3 | 電源 |
| GND | GND | グラウンド |
| DIN | GPIO2 | データ入力（1本のピンで35個を直列制御） |

配線を変更する場合は、電子ペーパー側は [epd.go](epd.go) 冒頭の `epd*Pin` 変数、円形ディスプレイ側は [round.go](round.go) 冒頭の `gc*Pin` 変数、
LEDマトリクス側は [matrix.go](matrix.go) 冒頭の `matrixLedPin` を書き換えること。

## 表示内容の変更

- [main.go](main.go) の `badgeName` / `badgeHandle` 定数で電子ペーパーの表示内容を変更できる（`badgeName` は [font.go](font.go) の `hiragana` マップに登録された文字のみ表示可能）。
- 円形ディスプレイのアイコンを変更する手順は [test11_dual_badge/README.md](../test11_dual_badge/README.md) を参照。
- LEDマトリクスの表示文字は [matrix.go](matrix.go) の `matrixMessage` を、フォントを増やす場合は `matrixGlyphs` にドットパターンを追加する。スクロール速度は `matrixScrollDelay`、明るさは `matrixBrightness`、色が変わる速さは `matrixHueStep` で調整できる。

## 必要なツール

- [TinyGo](https://tinygo.org/) (v0.41.1で動作確認)
- Go 1.25以降

## ビルド方法

```sh
cd test12_triple_badge
tinygo build -target=pico2 -o /tmp/out.uf2 .   # ビルド確認のみ
```

## 書き込み方法

```sh
tinygo flash -target=pico2 .   # BOOTSELモードで接続した状態で実行
```
