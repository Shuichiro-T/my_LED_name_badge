# test_06_GC901A_python

Raspberry Pi Zero + GC9A01A（1.28インチ 240x240 円形TFT、4線SPI）の動作確認用 Python プログラム。

CircuitPython（Adafruit Blinka）と Pillow を使い、`adafruit_rgb_display.gc9a01a` ドライバ経由で
PIL画像をそのまま表示する。

## 配線（標準SPI0ピン + 任意GPIO）

| GC9A01A | Raspberry Pi |
| ------- | ------------ |
| VCC     | 3.3V         |
| GND     | GND          |
| SCL     | GPIO11 (SCLK) |
| SDA     | GPIO10 (MOSI) |
| CS      | GPIO8 (CE0)  |
| DC      | GPIO25       |
| RST     | GPIO27       |

## セットアップ

```sh
sudo raspi-config  # Interface Options -> SPI を有効化してから再起動
pip3 install -r requirements.txt
sudo apt-get install -y fonts-dejavu
```

## 実行

```sh
python3 main.py
```

カラーバー → 十字線・円・テキスト → 単色塗り（赤・緑・青）の順に表示を切り替える。
