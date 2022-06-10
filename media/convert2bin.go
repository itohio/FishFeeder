package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	// "tinygo.org/x/drivers/image/png"
	"image/color"
	"image/png"
)

// See ../../image/README.md for the usage.

func main() {
	err := run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

// var buf [3 * 256]uint16

func RGBATo565(c color.Color) uint16 {
	r, g, b, _ := c.RGBA()
	return uint16((r & 0xF800) +
		((g & 0xFC00) >> 5) +
		((b & 0xF800) >> 11))
}

func run(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: %s FILE", args[0])
	}
	fname := args[1]
	varName := path.Base(fname)
	ext := path.Ext(fname)
	varName = strings.Title(strings.Replace(varName, ext, "", 1))
	ext = strings.Title(strings.Replace(ext, ".", "", 1))

	if ext != "Png" && ext != "Jpg" && ext != "Jpeg" {
		return fmt.Errorf("file %s is neither png nor a jpeg: %s", fname, ext)
	}

	b, err := os.Open(fname)
	if err != nil {
		return err
	}

	D := make([]uint16, 0, 40*40)
	var W, H int16
	// png.SetCallback(buf[:], func(data []uint16, x, y, w, h, width, height int16) {
	// 	D = append(D, data...)
	// 	W, H = width, height
	// })
	img, err := png.Decode(b)
	if err != nil {
		return err
	}
	W = int16(img.Bounds().Dx())
	H = int16(img.Bounds().Dy())
	for y := 0; y < int(H); y++ {
		for x := 0; x < int(W); x++ {
			c := img.At(x, y)
			D = append(D, RGBATo565(c))
		}
	}

	fmt.Println("package icons")
	fmt.Println()

	fmt.Println("const (")
	fmt.Printf("  %sWidth = %d\n", varName, W)
	fmt.Printf("  %sHeight = %d\n", varName, H)
	fmt.Println(")")
	fmt.Println("var (")
	fmt.Printf("  %s%s = []uint16 { \n", varName, ext)

	max := 10
	for i, d := range D {
		if (i % max) == 0 {
			fmt.Printf("    ")
		}

		fmt.Printf("%d, ", d)

		if (i%max) == 0 && i != 0 {
			fmt.Println()
		}
	}

	fmt.Println("  }")
	fmt.Println(")")

	return nil
}
