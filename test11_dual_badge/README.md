# test11_dual_badge

電子ペーパーモジュール（test07_EPD213BWR）と円形ディスプレイ（test02_GC9A01系）を1台のPico 2に同時接続し、それぞれ役割分担して表示する名札プログラムです。TinyGo で実装しています。

- 電子ペーパー（WeAct Studio 2.13インチ EPD213BWR, SSD1680, 122×250px, 白黒赤）: 名前「しゅういちろ」とXアカウント「@shucho0103」を表示。
- 円形ディスプレイ（GC9A01, 240×240px）: [icon/myicon.jpg](icon/myicon.jpg) を変換したアイコン画像を表示。

## 接続（配線）

2つのデバイスは独立したSPIバス・ピンを使うため、同時に接続しても干渉しません。ピンヘッダーに出ているGPIO0-15,26-29の範囲に収まり、かつRP2350のSPI0/SPI1が使えるGPIOの組み合わせ（決まった数パターンしかない）を満たすものを選んでいます。

### 電子ペーパー側（SPI1 + GPIO8-14）test07_EPD213BWRと同一の配線

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

### 円形ディスプレイ側（SPI0 + GPIO0,1,4-7,15）

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

配線を変更する場合は、電子ペーパー側は [epd.go](epd.go) 冒頭の `epd*Pin` 変数、円形ディスプレイ側は [round.go](round.go) 冒頭の `gc*Pin` 変数を書き換えること。

## 表示内容の変更

[main.go](main.go) の `badgeName` / `badgeHandle` 定数を書き換えれば、電子ペーパーの表示内容を変更できる。

- `badgeName` の文字は [font.go](font.go) の `hiragana` マップに登録されているグリフ（`し`,`ゅ`,`う`,`い`,`ち`,`ろ`）のみ表示可能。他の文字を使いたい場合はグリフデータを追加生成する必要がある（test07_EPD213BWR参照）。
- `badgeHandle` はASCII文字であれば任意の文字列を表示可能（`tinygo.org/x/tinyfont/freemono` フォントで描画）。

円形ディスプレイのアイコンを変更したい場合は、[icon/myicon.jpg](icon/myicon.jpg) を差し替えた上で、以下のPythonスクリプト（Pillowが必要: `pip install Pillow`）で [icon_bitmap.go](icon_bitmap.go) を再生成すること。

- 240x240pxにリサイズした上で、赤み(R高・G/B低)の画素を「インク」、それ以外を「背景」として2値化し、1bpp（1行30byte、MSBファースト）のビットマップとしてGoソースに埋め込んでいる。
- そのままフルカラー画像として埋め込むとGoソースが肥大化しすぎるため、この用途ではロゴ・印影のような単色マークに限られる。写真等フルカラー画像を表示したい場合は別途相談すること。

```python
from PIL import Image

SRC = "icon/myicon.jpg"
SIZE = 240

img = Image.open(SRC).convert("RGB").resize((SIZE, SIZE), Image.LANCZOS)

row_bytes = SIZE // 8
out_bytes = bytearray(row_bytes * SIZE)
for y in range(SIZE):
    for x in range(SIZE):
        r, g, b = img.getpixel((x, y))
        is_ink = r > 100 and g < r * 0.6 and b < r * 0.6
        if is_ink:
            out_bytes[y * row_bytes + x // 8] |= 0x80 >> (x % 8)

with open("icon_bitmap.go", "w") as f:
    f.write("package main\n\n")
    f.write(f"const xIconSize = {SIZE}\n\n")
    f.write("var xIconBitmap = [...]byte{\n")
    for y in range(SIZE):
        row = out_bytes[y * row_bytes:(y + 1) * row_bytes]
        f.write("\t" + ", ".join(f"0x{b:02x}" for b in row) + ",\n")
    f.write("}\n")
```

## 必要なツール

- [TinyGo](https://tinygo.org/) (v0.41.1で動作確認)
- Go 1.25以降

## ビルド方法

```sh
cd test11_dual_badge
tinygo build -target=pico2 -o /tmp/out.uf2 .   # ビルド確認のみ
```

## 書き込み方法

```sh
tinygo flash -target=pico2 .   # BOOTSELモードで接続した状態で実行
```
