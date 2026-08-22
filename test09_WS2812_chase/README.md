# test09_WS2812_chase

WS2812（NeoPixel互換）RGB LEDを35個直列に接続し、赤い光が先頭（1個目）から末尾（35個目）へ移動していくサンプルです。TinyGo で実装しています。

- `tinygo.org/x/drivers/ws2812` を利用（RP2040/RP2350のクロック周波数に応じたビットバンギング実装が含まれており、pico2にそのまま対応）。
- 1つのLEDだけを赤色に点灯させ、残りは消灯。点灯位置を先頭から末尾へ1つずつずらしながら繰り返すことで光が移動しているように見せている。末尾まで到達したら先頭に戻る。
- 明るさは眩しさ・消費電流対策で控えめ（`brightness = 0.1`）にしている。

## 接続（配線）

| WS2812テープ側 | Pico 2 側 | 役割 |
|---|---|---|
| VCC (5V or 3V3) | 5V または 3V3 | 電源 |
| GND | GND | グラウンド |
| DIN | GP15 | データ入力（1本のピンで35個を直列制御） |

配線を変更する場合は [main.go](main.go) 冒頭の `ledPin` を書き換えること。LED数や明るさを変える場合は `numLEDs` / `brightness` を変更する。

## ビルド方法

```sh
cd test09_WS2812_chase
tinygo build -target=pico2 -o /tmp/out.uf2 .   # ビルド確認のみ
```

## 書き込み方法

```sh
tinygo flash -target=pico2 .   # BOOTSELモードで接続した状態で実行
```
