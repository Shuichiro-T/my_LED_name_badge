# test08_WS2812_rainbow

WS2812（NeoPixel互換）RGB LEDを9個直列に接続し、虹色のグラデーションが流れるように点灯させるサンプルです。TinyGo で実装しています。

- `tinygo.org/x/drivers/ws2812` を利用（RP2040/RP2350のクロック周波数に応じたビットバンギング実装が含まれており、pico2にそのまま対応）。
- HSV色空間で各LEDの色相をずらして虹のグラデーションを作り、時間経過で色相全体を回転させることでアニメーションさせている。

## 接続（配線）

| WS2812テープ側 | Pico 2 側 | 役割 |
|---|---|---|
| VCC (5V or 3V3) | 5V または 3V3 | 電源 |
| GND | GND | グラウンド |
| DIN | GP15 | データ入力（1本のピンで9個を直列制御） |

配線を変更する場合は [main.go](main.go) 冒頭の `ledPin` を書き換えること。LED数を変える場合は `numLEDs` を変更する。

## ビルド方法

```sh
cd test08_WS2812_rainbow
tinygo build -target=pico2 -o /tmp/out.uf2 .   # ビルド確認のみ
```

## 書き込み方法

```sh
tinygo flash -target=pico2 .   # BOOTSELモードで接続した状態で実行
```
