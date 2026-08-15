# test02_GC9A01

RP2350-Zero (Waveshare) + GC9A01 円形ディスプレイで動作するアナログ時計です。TinyGo で実装しています。

## 接続（配線）

RP2350-Zeroには`pico2`のようなボード固有の名前付きピン（GP0など）が無いため、RP2350のGPIO番号をそのまま使用します。

RP2350-Zeroは29本のGPIOのうち、**GP17〜GP25の9本は基板裏面の半田付けが必要なパッド**で、それ以外（GP0〜GP15, GP26〜GP29の計20本）はエッジピンアウト（ヘッダピン）から半田付け無しで利用できます。本プログラムではGP17〜GP25を避け、エッジピンアウトのみを使用しています。

SPI0と制御用GPIOを以下のように接続してください。

| GC9A01 モジュール側 | RP2350-Zero 側 | 役割 |
|---|---|---|
| VCC | 3V3 (OUT) | 電源 (3.3V) |
| GND | GND | グラウンド |
| SCL / SCK | GPIO2 | SPIクロック |
| SDA / MOSI | GPIO3 | SPIデータ (Tx) |
| RES / RESET | GPIO4 | リセット |
| DC | GPIO5 | データ/コマンド切替 |
| CS | GPIO1 | チップセレクト |
| BLK / BL | GPIO9 | バックライト |

- GC9A01モジュールにMISO(SDO)端子がある場合はGPIO0に接続できますが、このプログラムでは読み取りを行わないため未接続でも構いません。
- BLK端子が無いモジュール（常時バックライトON）の場合はGPIO9への接続は不要です。
- オンボードのWS2812 RGB LEDはGPIO16（半田パッド側）を使用しています。
- 配線を変更したい場合は [main.go](main.go) 内の以下の箇所を書き換えてください（GP17〜GP25は避けること）。

```go
// RP2350-Zero (Waveshare) のピン配置。GP17-GP25は半田付けが必要なため、
// エッジピンアウトのみ（GP0-GP15, GP26-GP29）から選んでいる。
machine.SPI0.Configure(machine.SPIConfig{
    Frequency: 80000000,
    SCK:       machine.GPIO2,
    SDO:       machine.GPIO3, // MOSI
    SDI:       machine.GPIO0, // MISO (未使用)
})

display := gc9a01.New(machine.SPI0, machine.GPIO4, machine.GPIO5, machine.GPIO1, machine.GPIO9)
// 引数の順番: SPIバス, リセットピン, DCピン, CSピン, バックライトピン
```

## 時刻の設定について

RP2350には電池バックアップ付きのRTCが無いため、起動時刻は [main.go](main.go) 冒頭の `initialTime` に埋め込んでいます。書き込み（フラッシュ）する直前の時刻に書き換えてからビルドしてください。

```go
var initialTime = time.Date(2026, 6, 20, 12, 0, 0, 0, time.UTC)
```

以降は内部タイマーでこの時刻からの経過時間を加算して表示します（USBケーブルを挿し直すなど再起動すると、また `initialTime` の時刻からスタートします）。

## 必要なツール

- [TinyGo](https://tinygo.org/) (このプログラムは v0.40.1 で動作確認)
- Go 1.21以降

## ビルド方法

RP2350-Zero専用のTinyGoボード定義は無いため、リポジトリルートに置いた [rp2350-zero.json](../rp2350-zero.json)（カスタムターゲット定義）を使用します。

> [!NOTE]
> 汎用の `rp2350` ターゲットは `inheritable-only`（他ターゲットに継承されるための基底定義）に指定されているため、TinyGo 0.41系以降では `-target=rp2350` を直接指定するとビルド・書き込みの両方で `target is inheritable-only, which means it cannot be used directly for building or flashing` というエラーになります。そのため Seeed XIAO RP2350 向けのボード定義（`xiao-rp2350`。ピン別名は未使用なので実害なし）を継承しつつ、RP2350-Zeroの実機に合わせてフラッシュ容量を4MBに上書きした `rp2350-zero.json` を用意しています。

```sh
cd test02_GC9A01
tinygo build -target=../rp2350-zero.json -o test02_GC9A01.uf2 .
```

## 書き込み方法（USBブートローダ経由）

1. RP2350-Zero基板上の **BOOT** ボタンを押しながらUSBケーブルをPCに接続する（ボタンを押したまま接続し、接続後に離す）。
2. PCに `RP2350` という名前のUSBマスストレージドライブが現れる。
3. 生成された `test02_GC9A01.uf2` をそのドライブにコピー（ドラッグ&ドロップ）する。
4. コピー完了後、ボードが自動的に再起動し、プログラムが実行される。

`tinygo` コマンド一発で書き込みまで行いたい場合は、BOOTモードでPCに接続した状態で以下を実行してください。

```sh
tinygo flash -target=../rp2350-zero.json .
```

## 表示内容

- 240x240の円形画面いっぱいにアナログ時計の文字盤を描画します。
- 12時・3時・6時・9時の位置に数字、5分ごとに長め、1分ごとに短めの目盛りを表示します。
- 時針（白・短）・分針（白・長）・秒針（赤）を1秒ごとに再描画します。
