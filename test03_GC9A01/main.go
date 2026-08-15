package main

import (
	"machine"
	"time"

	"image/color"

	"tinygo.org/x/drivers/gc9a01"
)

func main() {
	machine.SPI1.Configure(machine.SPIConfig{
		Frequency: 80_000_000,
		SCK:       machine.GPIO10,
		SDO:       machine.GPIO11, // MOSI
		SDI:       machine.GPIO12, // MISO (未使用)
		Mode:      0,
	})

	display := gc9a01.New(machine.SPI1, machine.GPIO4, machine.GPIO5, machine.GPIO1, machine.GPIO9)
	display.Configure(gc9a01.Config{Orientation: gc9a01.HORIZONTAL, Width: 240, Height: 240})

	// width, height := display.Size()

	// white := color.RGBA{255, 255, 255, 255}
	red := color.RGBA{255, 0, 0, 255}
	// blue := color.RGBA{0, 0, 255, 255}
	//green := color.RGBA{0, 255, 0, 255}
	// black := color.RGBA{0, 0, 0, 255}

	display.FillScreen(red)

	// display.FillRectangle(0, 0, width/2, height/2, white)
	// display.FillRectangle(width/2, 0, width/2, height/2, red)
	// display.FillRectangle(0, height/2, width/2, height/2, green)
	// display.FillRectangle(width/2, height/2, width/2, height/2, blue)
	// display.FillRectangle(width/4, height/4, width/2, height/2, black)

	for {
		time.Sleep(time.Hour)
	}

}
