# my_LED_name_badge

電子工作（電光式の名札／LEDバッジ）プロジェクト。マイコンボードには Raspberry Pi Pico 2 (RP2350) 系を使用し、プログラミングは主に TinyGo で行う。

`testNN_<デバイス名>/` という連番のサブディレクトリごとに、動作検証を積み重ねていく方針。各サブディレクトリは個別の Go module（`go.mod`）を持つ。

## ディレクトリ一覧

| ディレクトリ | 内容 |
|---|---|
| [test01_GC9A01](test01_GC9A01) | GC9A01円形ディスプレイの画面全体を赤色に点灯させるだけの最小構成の動作確認（TinyGo）。SPI通信の不安定さを切り分けるため、SPI1・別ピン構成を試行中。 |
| [test02_GC9A01](test02_GC9A01) | GC9A01円形ディスプレイにアナログ時計を描画するサンプル（TinyGo）。`tinydraw` / `tinyfont` を使って文字盤・針を描画する。 |
| [test03_GC9A01](test03_GC9A01) | GC9A01を`tinygo.org/x/drivers/gc9a01`ドライバで初期化し画面を赤色で塗りつぶす検証（TinyGo）。SPI1 + GPIO10-14の配線で導通確認済み。 |
| [test05_GC9A01_raw](test05_GC9A01_raw) | ドライバを使わず、GC9A01のレジスタ制御を直接SPIコマンドとして送信する低レイヤーの初期化・描画検証（TinyGo）。test03で確認した配線を踏襲。 |
| [test06_GC901A_python](test06_GC901A_python) | Raspberry Pi Zero + GC9A01Aを、CircuitPython（Adafruit Blinka）とPillowで動作確認するPythonプログラム。 |
| [test07_EPD213BWR](test07_EPD213BWR) | WeAct Studio 2.13インチ電子ペーパーモジュール（E0213A179、SSD1680、122x250px、白黒赤3色）に「しゅういちろ」と表示する名札プログラム（TinyGo）。ドライバはSSD1680を直接SPI制御する自前実装。 |

## 対象ハードウェア

- マイコンボード: Raspberry Pi Pico 2 (RP2350) / RP2350-Zero (Waveshare)
- TinyGo ビルドターゲット: `pico2`（RP2350-Zeroの場合は[rp2350-zero.json](rp2350-zero.json)を使用）
