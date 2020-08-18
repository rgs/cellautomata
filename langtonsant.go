package main

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"os"
)

var palette []color.Color
const paletteSize = 256
const width, height = 512, 512

// Describes rules for each color
var turns []func(int)int

const (
    up = iota
    rt
    dn
    lt
)

func turnRight(dir int) (d int) {
	switch dir {
	case up: d = rt
	case rt: d = dn
	case dn: d = lt
	case lt: d = up
	}
	return
}

func turnLeft(dir int) (d int) {
	switch dir {
	case up: d = lt
	case rt: d = up
	case dn: d = rt
	case lt: d = dn
	}
	return
}

func turnBack(dir int) (d int) {
	switch dir {
	case up: d = dn
	case rt: d = lt
	case dn: d = up
	case lt: d = rt
	}
	return
}

func turnNone(dir int) int { return dir }

func init() {
	palette = make([]color.Color, paletteSize)
	turns = make([]func(int)int, paletteSize)
	for i := 0 ; i < paletteSize ; i++ {
		g := 0xff - uint8(i * 0xff / (paletteSize - 1))
		palette[i] = color.NRGBA{g, g, g, 0xff}
		// This is going to generate a RLLRRLL...R pattern
		switch ((i + 1) / 2) % 2 {
		case 0: turns[i] = turnRight
		case 1: turns[i] = turnLeft
		}
	}
}

func findColorInPalette(c color.Color) int {
	// TODO a lookup table would be nice...
	for i := 0 ; i < paletteSize ; i++ {
		if palette[i] == c {
			return i
		}
	}
	log.Fatalln("color not found", c)
	return 0
}

func main() {
	bounds := image.Rect(0, 0, width, height)
	img := image.NewNRGBA(bounds)
	draw.Draw(img, bounds, image.NewUniform(palette[0]), image.ZP, draw.Src)
	dir := lt
	pos := image.Point{width / 2, height / 2}
	iter := 0
	log.Println("Starting")
	for pos.In(bounds) {
		curColor := img.At(pos.X, pos.Y).(color.NRGBA)
		i := findColorInPalette(curColor)
		img.Set(pos.X, pos.Y, palette[(i + 1) % paletteSize])
		dir = turns[i](dir)
		switch dir {
		case up: pos.Y -= 1
		case rt: pos.X += 1
		case dn: pos.Y += 1
		case lt: pos.X -= 1
		}
		iter++
	}
	log.Println("Iterations:", iter)
	f, err := os.Create("langtonsant.png")
	if err != nil {
		log.Fatalln(err)
	}
	if err := png.Encode(f, img); err != nil {
		f.Close()
		log.Fatalln(err)
	}
	if err := f.Close(); err != nil {
		log.Fatalln(err)
	}
}
