# test01_GC9A01

RP2350-Zero (Waveshare) + GC9A01 円形ディスプレイの画面全体を赤色に点灯させるだけの動作確認用プログラムです。TinyGo で実装しています。

> [!NOTE]
> 現在、通信不安定の切り分けのためSPI0からSPI1・別ピン構成に一時的に変更中です。そのため [test02_GC9A01](../test02_GC9A01) とはピン割り当てが異なります。問題解決後は元のSPI0構成に戻す予定です。

## 接続（配線）

RP2350-Zeroには`pico2`のようなボード固有の名前付きピン（GP0など）が無いため、RP2350のGPIO番号をそのまま使用します。

RP2350-Zeroは29本のGPIOのうち、**GP17〜GP25の9本は基板裏面の半田付けが必要なパッド**で、それ以外（GP0〜GP15, GP26〜GP29の計20本）はエッジピンアウト（ヘッダピン）から半田付け無しで利用できます。本プログラムではGP17〜GP25を避け、エッジピンアウトのみを使用しています。

SPI0と制御用GPIOを以下のように接続してください。

| GC9A01 モジュール側 | RP2350-Zero 側 | 役割 |
|---|---|---|
| VCC | 3V3 (OUT) | 電源 (3.3V) |
| GND | GND | グラウンド |
| SCL / SCK | GPIO10 | SPIクロック |
| SDA / MOSI | GPIO11 | SPIデータ (Tx) |
| RES / RESET | GPIO12 | リセット |
| DC | GPIO13 | データ/コマンド切替 |
| CS | GPIO9 | チップセレクト（常時Low固定） |
| BLK / BL | GPIO14 | バックライト |

- GC9A01モジュールにMISO(SDO)端子がある場合はGPIO8に接続できますが、このプログラムでは読み取りを行わないため未接続でも構いません。
- CSピンはドライバが描画時に明示的に操作しないため、起動時に一度Lowへ固定し、常時選択状態のまま使用しています。
- オンボードのWS2812 RGB LEDはGPIO16（半田パッド側）を使用しています。
- 配線を変更したい場合は [main.go](main.go) 内の以下の箇所を書き換えてください（GP17〜GP25は避けること）。

```go
// [切り分け用に一時変更中] SPI1 + 別ピンに配線を移して同じ症状が出るか確認する。
machine.SPI1.Configure(machine.SPIConfig{
    Frequency: 1000000, // さらに低速にして信号品質要因を排除
    SCK:       machine.GPIO10,
    SDO:       machine.GPIO11, // MOSI
    SDI:       machine.GPIO8,  // MISO (未使用)
})

// CSピンは常時選択状態（Low）に固定する。
csPin := machine.GPIO9
csPin.Configure(machine.PinConfig{Mode: machine.PinOutput})
csPin.Low()

display := gc9a01.New(machine.SPI1, machine.GPIO12, machine.GPIO13, csPin, machine.GPIO14)
// 引数の順番: SPIバス, リセットピン, DCピン, CSピン, バックライトピン
```

## 必要なツール

- [TinyGo](https://tinygo.org/) (このプログラムは v0.40.1 で動作確認)
- Go 1.21以降

## ビルド方法

RP2350-Zero専用のTinyGoボード定義は無いため、リポジトリルートに置いた [rp2350-zero.json](../rp2350-zero.json)（カスタムターゲット定義）を使用します。

> [!NOTE]
> 汎用の `rp2350` ターゲットは `inheritable-only`（他ターゲットに継承されるための基底定義）に指定されているため、TinyGo 0.41系以降では `-target=rp2350` を直接指定するとビルド・書き込みの両方で `target is inheritable-only, which means it cannot be used directly for building or flashing` というエラーになります。そのため Seeed XIAO RP2350 向けのボード定義（`xiao-rp2350`。ピン別名は未使用なので実害なし）を継承しつつ、RP2350-Zeroの実機に合わせてフラッシュ容量を4MBに上書きした `rp2350-zero.json` を用意しています。

```sh
cd test01_GC9A01
tinygo build -target=../rp2350-zero.json -o test01_GC9A01.uf2 .
```

## 書き込み方法（USBブートローダ経由）

1. RP2350-Zero基板上の **BOOT** ボタンを押しながらUSBケーブルをPCに接続する（ボタンを押したまま接続し、接続後に離す）。
2. PCに `RP2350` という名前のUSBマスストレージドライブが現れる。
3. 生成された `test01_GC9A01.uf2` をそのドライブにコピー（ドラッグ&ドロップ）する。
4. コピー完了後、ボードが自動的に再起動し、プログラムが実行される。

`tinygo` コマンド一発で書き込みまで行いたい場合は、BOOTモードでPCに接続した状態で以下を実行してください。

```sh
tinygo flash -target=../rp2350-zero.json .
```

## 表示内容

- 画面全体（240x240）を赤色一色で塗りつぶして表示し続けます。配線確認用の最小構成です。
