"""
Raspberry Pi Zero + GC9A01A (1.28インチ 240x240 円形TFT) 動作確認用テストパターン

配線（標準SPI0ピン + 任意GPIO）:
    VCC -> 3.3V
    GND -> GND
    SCL -> GPIO11 (SCLK)
    SDA -> GPIO10 (MOSI)
    CS  -> GPIO8  (CE0)
    DC  -> GPIO25
    RST -> GPIO27

事前準備:
    sudo raspi-config  # Interface Options -> SPI を有効化
    pip3 install -r requirements.txt
    sudo apt-get install -y fonts-dejavu

実行:
    python3 main.py
"""

import time

import board
import digitalio
from PIL import Image, ImageDraw, ImageFont
from adafruit_rgb_display import gc9a01a

WIDTH = 240
HEIGHT = 240
# 縞模様が出る場合はSPIクロックが高すぎてジャンパ線でノイズが乗っている可能性があるため、
# まず低速（1MHz）で動作確認する。安定したら徐々に上げる。
BAUDRATE = 1_000_000

FONT_PATH = "/usr/share/fonts/truetype/dejavu/DejaVuSans-Bold.ttf"


def create_display():
    spi = board.SPI()
    cs_pin = digitalio.DigitalInOut(board.CE0)
    dc_pin = digitalio.DigitalInOut(board.D25)
    rst_pin = digitalio.DigitalInOut(board.D27)

    return gc9a01a.GC9A01A(
        spi,
        cs=cs_pin,
        dc=dc_pin,
        rst=rst_pin,
        width=WIDTH,
        height=HEIGHT,
        baudrate=BAUDRATE,
    )


def draw_color_bars(draw):
    colors = [
        (255, 0, 0),
        (0, 255, 0),
        (0, 0, 255),
        (255, 255, 0),
        (0, 255, 255),
        (255, 0, 255),
        (255, 255, 255),
        (0, 0, 0),
    ]
    bar_height = HEIGHT // len(colors)
    for i, color in enumerate(colors):
        y0 = i * bar_height
        y1 = HEIGHT if i == len(colors) - 1 else y0 + bar_height
        draw.rectangle((0, y0, WIDTH, y1), fill=color)


def draw_crosshair_and_circle(draw):
    draw.line((WIDTH // 2, 0, WIDTH // 2, HEIGHT), fill=(0, 0, 0), width=2)
    draw.line((0, HEIGHT // 2, WIDTH, HEIGHT // 2), fill=(0, 0, 0), width=2)
    margin = 4
    draw.ellipse((margin, margin, WIDTH - margin, HEIGHT - margin), outline=(0, 0, 0), width=3)


def draw_text_overlay(image, draw):
    font = ImageFont.truetype(FONT_PATH, 28)
    text = "GC9A01A TEST"
    bbox = draw.textbbox((0, 0), text, font=font)
    text_w = bbox[2] - bbox[0]
    text_h = bbox[3] - bbox[1]
    draw.text(
        ((WIDTH - text_w) // 2, (HEIGHT - text_h) // 2),
        text,
        font=font,
        fill=(0, 0, 0),
    )


def main():
    disp = create_display()

    # 1. カラーバー
    image = Image.new("RGB", (WIDTH, HEIGHT))
    draw = ImageDraw.Draw(image)
    draw_color_bars(draw)
    disp.image(image)
    time.sleep(3)

    # 2. 十字線・円・テキスト
    image = Image.new("RGB", (WIDTH, HEIGHT), (255, 255, 255))
    draw = ImageDraw.Draw(image)
    draw_crosshair_and_circle(draw)
    draw_text_overlay(image, draw)
    disp.image(image)
    time.sleep(3)

    # 3. 単色塗り（赤 -> 緑 -> 青）を繰り返し表示
    for color in [(255, 0, 0), (0, 255, 0), (0, 0, 255)]:
        image = Image.new("RGB", (WIDTH, HEIGHT), color)
        disp.image(image)
        time.sleep(1)


if __name__ == "__main__":
    main()
