package main

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"runtime/pprof"
	"strings"
)

func ycost(img image.Image, stride, offset, best int) (cost int) {
	X, Y := img.Bounds().Size().X, img.Bounds().Size().Y
	fails := 0
	for x := 0; x < X; x++ {
		failing := false
		for y := stride + offset; y < Y; y += stride {
			r, g, b, _ := img.At(x, y).RGBA()
			if r != 0 || g != 0 || b != 0 {
				failing = true
			} else if failing {
				failing = false
				fails++

				if fails > best {
					return fails
				}
			}
		}
	}
	return fails
}

func xcost(img image.Image, stride, offset, best int) (cost int) {
	X, Y := img.Bounds().Size().X, img.Bounds().Size().Y
	fails := 0
	for y := 0; y < Y; y++ {
		failing := false
		for x := stride + offset; x < X; x += stride {
			r, g, b, _ := img.At(x, y).RGBA()
			if r != 0 || g != 0 || b != 0 {
				failing = true
			} else if failing {
				failing = false
				fails++
				if fails > best {
					return fails
				}
			}
		}
	}
	return fails
}

/*
func xpaint(img *image.RGBA, stride int) {
	for y := 0; y < img.Bounds().Size().Y; y++ {
		for x := stride - 1; x < img.Bounds().Size().X; x += stride {
			img.Set(x, y, color.RGBA{0, 0, 255, 255})
		}
	}
}
*/

type textcons struct {
	lut map[string]rune
}

func readPNG(filename string) (image.Image, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	// TODO: Make errcheck happy.  Errcheck should probably ignore
	// errors when closing read only files
	defer func() {
		if err := f.Close(); err != nil {
			log.Println(err)
		}
	}()

	return png.Decode(f)
}

func calculateStrides(img image.Image) (int, int, int, int, error) {
	var ystride, xstride int
	var yoff, xoff int
	ybest, xbest := 9999999, 9999999

	//fmt.Println("max Y stride", img.Bounds().Size().Y/12)
outery:
	for stride := 2; stride < 17; stride++ {
		for off := -((stride - 1) / 2); off < stride-1; off++ {
			cost := ycost(img, stride, off, ybest)
			//fmt.Println("y stride", stride, "off", off, "cost", cost)
			if cost < ybest {
				ystride = stride
				yoff = off
				ybest = cost
				if ybest == 0 {
					break outery
				}
			}
		}
	}

outerx:
	for stride := 2; stride < 15; stride++ {
		for off := -((stride - 1) / 2); off < stride-1; off++ {
			cost := xcost(img, stride, off, xbest)
			//fmt.Println("x stride", stride, "off", off, "cost", cost)
			if cost < xbest {
				xstride = stride
				xoff = off
				xbest = cost
				if xbest == 0 {
					break outerx
				}
			}
		}
	}
	fmt.Println(xstride, ystride, xoff, yoff, xbest, ybest)
	//	fmt.Println("ycost", ycost(img, ystride, yoff))
	//	fmt.Println("xcost", xcost(img, xstride, xoff))
	var err error
	if xstride == 0 || ystride == 0 {
		err = errors.New("stride unknown")
	}

	return xstride, ystride, xoff, yoff, err
}

func debugStrides(img image.Image, xstride, ystride, xoff, yoff int) {
	cost := 0
	out := image.NewRGBA(img.Bounds())
	for x := 0; x < img.Bounds().Size().X; x++ {
		for y := 0; y < img.Bounds().Size().Y; y++ {
			if ((y-yoff)%ystride == 0) || ((x-xoff)%xstride == 0) {
				r, _, _, _ := img.At(x, y).RGBA()
				if r != 0 {
					out.Set(x, y, color.RGBA{0, 255, 0, 255})
					cost++
				} else {
					out.Set(x, y, color.RGBA{255, 0, 0, 255})
				}
			} else {
				out.Set(x, y, img.At(x, y))
			}
		}
	}

	fmt.Println("cost", cost)

	fout, err := os.Create("out.png")
	if err != nil {
		log.Fatal(err)
	}
	if err := png.Encode(fout, out); err != nil {
		log.Fatal(err)
	}
	if err := fout.Close(); err != nil {
		log.Fatal(err)
	}

}

func iterateCells(img image.Image, xstride, ystride, xoff, yoff int, cb func(col, row int, pixmap []byte)) {
	X, Y := img.Bounds().Size().X, img.Bounds().Size().Y
	bytes := make([]byte, xstride*ystride)
	var row int
	for y := yoff; y+ystride < Y; y += ystride {
		var col int
		for x := xoff; x+xstride < X; x += xstride {
			i := 0
			for dx := 0; dx < xstride; dx++ {
				for dy := 0; dy < ystride; dy++ {
					r, g, b, _ := img.At(x+dx, y+dy).RGBA()
					if r != 0 || g != 0 || b != 0 {
						bytes[i] = 255
					} else {
						bytes[i] = 0
					}
					i++
				}
			}
			cb(col, row, bytes)
			col++
		}
		row++
	}
}

func (t *textcons) learn(imgfile, textfile string) error {
	if t.lut == nil {
		t.lut = make(map[string]rune)
	}

	img, err := readPNG(imgfile)
	if err != nil {
		return err
	}

	xstride, ystride, xoff, yoff, err := calculateStrides(img)
	if err != nil {
		return err
	}

	text, err := ioutil.ReadFile(textfile)
	if err != nil {
		return err
	}

	var lines [][]rune
	for _, line := range strings.Split(string(text), "\n") {
		lines = append(lines, []rune(line))
	}

	isEmpty := func(pixmap []byte) bool {
		for _, b := range pixmap {
			if b != 0 {
				return false
			}
		}
		return true
	}

	iterateCells(img, xstride, ystride, xoff, yoff, func(col, row int, pixmap []byte) {
		key := make([]byte, len(pixmap))
		for i, b := range pixmap {
			if b != 0 {
				key[i] = '1'
			} else {
				key[i] = '0'
			}
		}
		if isEmpty(pixmap) {
			t.lut[string(key)] = ' '
		}

		if row >= len(lines) || col >= len(lines[row]) {
			if !isEmpty(pixmap) {
				fmt.Println("missing text for", col, row, string(key))
			}

			return
		}

		if r, ok := t.lut[string(key)]; ok {
			if r != lines[row][col] {
				fmt.Println("already defined", col, row, string(key), string([]rune{lines[row][col], r}))
			}
			return
		}
		t.lut[string(key)] = lines[row][col]
		//fmt.Println(col, row, string(key), string([]rune{lines[row][col]}))
	})

	return nil
}

func (t *textcons) parse(imgfile string) ([]string, error) {
	img, err := readPNG(imgfile)
	if err != nil {
		return nil, err
	}

	xstride, ystride, xoff, yoff, err := calculateStrides(img)
	if err != nil {
		return nil, err
	}

	debugStrides(img, xstride, ystride, xoff, yoff)

	var lines []*bytes.Buffer

	iterateCells(img, xstride, ystride, xoff, yoff, func(col, row int, pixmap []byte) {
		if row >= len(lines) {
			lines = append(lines, &bytes.Buffer{})
		}

		key := make([]byte, len(pixmap))
		for i, b := range pixmap {
			if b != 0 {
				key[i] = '1'
			} else {
				key[i] = '0'
			}
		}

		if r, ok := t.lut[string(key)]; ok {
			fmt.Fprintf(lines[row], "%c", r)
		} else {
			fmt.Fprintf(lines[row], "?")
		}
	})
	var ret []string
	for _, row := range lines {
		ret = append(ret, row.String())
	}
	return ret, nil
}

func main() {
	{
		f, err := os.Create("prof")
		if err != nil {
			log.Fatal(err)
		}
		err = pprof.StartCPUProfile(f)
		if err != nil {
			log.Fatal(err)
		}
		defer pprof.StopCPUProfile()
	}

	var tc textcons
	err := tc.learn("testdata/abc.png", "testdata/abc.text")
	if err != nil {
		log.Fatal(err)
	}

	err = tc.learn("testdata/Installer_01_launch_early.png", "testdata/Installer_01_launch_early.text")
	if err != nil {
		log.Fatal(err)
	}

	err = tc.learn("testdata/Codepage-437.png", "testdata/Codepage-437.text")
	if err != nil {
		log.Fatal(err)
	}

	err = tc.learn("testdata/vga8x16.png", "testdata/vga8x16.text")
	if err != nil {
		log.Fatal(err)
	}

	text, err := tc.parse(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	for _, line := range text {
		fmt.Println(line)
	}

}
