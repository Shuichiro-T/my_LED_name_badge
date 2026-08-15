# my_LED_name_badge

電子工作（電光式の名札／LEDバッジ）プロジェクト。マイコンのプログラミングは TinyGo で行う。

## 出力ルール

- Claude からの応答・コメント・コミットメッセージなどの出力はすべて日本語で行う。

## 対象ハードウェア

- マイコンボード: Raspberry Pi Pico 2 (RP2350)
- TinyGo ビルドターゲット: `pico2`
- 接続確認済みシリアルポート例: `/dev/cu.usbmodem11301`（環境により変わる）

## ディレクトリ構成

- `testNN_<デバイス名>/` のように、検証用のサブディレクトリを連番で作成していく方針（例: `test02_GC9A01/` は GC9A01 円形LCDの動作検証）。
- 各サブディレクトリは個別の Go module（`go.mod`）を持つ。

## ビルド・書き込みコマンド

```sh
cd test02_GC9A01
tinygo build -target=pico2 -o /tmp/out.uf2 .   # ビルド確認のみ
tinygo flash -target=pico2 .                  # 実機への書き込み（BOOTSELモードで接続）
```

## 既存実装メモ

- `test02_GC9A01/main.go`: GC9A01（240x240, SPI接続）に時計（アナログ時計表示）を描画するサンプル。`machine.SPI0` と `GP6`–`GP9` を使用。`tinygo.org/x/drivers`, `tinydraw`, `tinyfont` を利用。
