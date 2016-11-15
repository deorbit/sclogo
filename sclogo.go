package main

import (
	"bufio"
	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dimg"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

type SpreadingCode struct {
	Chips []int
}

// Sequency returns the number of sign changes in the
// spreading code. Useful for creating sequency-ordered Walsh
// matrices.
func (code *SpreadingCode) Sequency() int {
	count := 0
	chips := code.Chips
	if len(chips) == 0 {
		return 0
	}

	for i := 1; i < len(chips); i++ {
		if chips[i] != chips[i-1] {
			count++
		}
	}

	return count
}

// WalshMatrix holds orthogonal spreading codes.
type WalshMatrix struct {
	Rows []SpreadingCode
}

func (m WalshMatrix) Len() int {
	return len(m.Rows)
}
func (m WalshMatrix) Swap(i, j int) {
	m.Rows[i], m.Rows[j] = m.Rows[j], m.Rows[i]
}
func (m WalshMatrix) Less(i, j int) bool {
	return m.Rows[i].Sequency() < m.Rows[j].Sequency()
}

// WalshMatrixFromFile reads a file containing a grid of positive and negative
// numbers representing chips in a Walsh matrix of orthogonal spreading codes.
func WalshMatrixFromFile(filename string) *WalshMatrix {
	matrixFile, err := os.Open(filename)
	matrix := WalshMatrix{}

	if err != nil {
		log.Fatal("Couldn't open matrix file.")
	}

	scanner := bufio.NewScanner(matrixFile)
	for scanner.Scan() {
		row := strings.Fields(scanner.Text())
		chips := make([]int, len(row))
		for i, chipStr := range row {
			chips[i], err = strconv.Atoi(chipStr)
			if err != nil {
				log.Printf("Couldn't convert %s to int.", chipStr)
			}
		}
		code := SpreadingCode{Chips: chips}
		matrix.Rows = append(matrix.Rows, code)
	}

	return &matrix
}

type Image struct {
	RGBA *image.RGBA
}

func NewImage() *Image {
	rgba := image.NewRGBA(image.Rect(0, 0, 100, 100))
	return &Image{RGBA: rgba}
}

// HLine draws a horizontal line
func (img *Image) HLine(x1, y, x2 int, color color.Color) {
	for ; x1 <= x2; x1++ {
		img.RGBA.Set(x1, y, color)
	}
}

// VLine draws a veritcal line
func (img *Image) VLine(x, y1, y2 int, color color.Color) {
	for ; y1 <= y2; y1++ {
		img.RGBA.Set(x, y1, color)
	}
}

// Rect draws a rectangle utilizing HLine() and VLine()
func (img *Image) Rect(x1, y1, x2, y2 int, color color.Color) {
	img.HLine(x1, y1, x2, color)
	img.HLine(x1, y2, x2, color)
	img.VLine(x1, y1, y2, color)
	img.VLine(x2, y1, y2, color)
}

func Rectangle(gc draw2d.GraphicContext, x1, y1, x2, y2 float64) {
	gc.MoveTo(x1, y1)
	gc.LineTo(x2, y1)
	gc.LineTo(x2, y2)
	gc.LineTo(x1, y2)
	gc.Close()
	gc.FillStroke()
}

func main() {
	matrix := WalshMatrixFromFile("walsh.txt")
	sort.Sort(matrix)

	logoSize := 128

	// Initialize the graphic context on an RGBA image
	dest := image.NewRGBA(image.Rect(0, 0, logoSize, logoSize))
	gc := draw2dimg.NewGraphicContext(dest)

	// Prepare the canvas.
	gc.SetLineWidth(0)

	matrixSize := len(matrix.Rows)
	cellPixSize := float64(logoSize / matrixSize)

	for i, row := range matrix.Rows {
		for j, col := range row.Chips {
			if col < 0 {
				gc.SetFillColor(color.RGBA{0x44, 0xff, 0x44, 0xff})
			} else {
				gc.SetFillColor(color.RGBA{0xff, 0x00, 0x44, 0x00})
			}
			Rectangle(gc,
				float64(i)*cellPixSize,
				float64(j)*cellPixSize,
				float64(i)*cellPixSize+cellPixSize,
				float64(j)*cellPixSize+cellPixSize)
		}
	}

	draw2dimg.SaveToPngFile("logo1.png", dest)

	var img = NewImage()
	var col color.Color

	col = color.RGBA{255, 0, 0, 255} // Red
	img.HLine(10, 20, 80, col)
	col = color.RGBA{0, 255, 0, 255} // Green
	img.Rect(10, 10, 80, 50, col)

	f, err := os.Create("logo.png")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	png.Encode(f, img.RGBA)
}
