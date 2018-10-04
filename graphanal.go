package main

import (
	"flag"
	"image/color"
	"log"
	"math"

	"github.com/fogleman/gg"
)

// Point represents a 2d point in spaces
type Point struct {
	x float64
	y float64
}

func newPoint(x int, y int) *Point {
	return &Point{float64(x), float64(y)}
}

func fillRecurse(width int, height int, context *gg.Context, startingColor color.Color, fillColor color.Color) {
	if width < 0 || height < 0 || width >= context.Width() || height >= context.Height() {
		return
	}
	currR, currG, currB, currA := context.Image().At(width, height).RGBA()
	startingR, startingG, startingB, startingA := startingColor.RGBA()
	if currR == startingR && currG == startingG && currB == startingB && currA == startingA {
		context.SetColor(fillColor)
		context.SetPixel(width, height)
		fillRecurse(width-1, height, context, startingColor, fillColor)
		fillRecurse(width+1, height, context, startingColor, fillColor)
		fillRecurse(width, height-1, context, startingColor, fillColor)
		fillRecurse(width, height+1, context, startingColor, fillColor)
	}
}

func fill(startWidth int, startHeight int, context *gg.Context, fillColor color.Color) {
	fillRecurse(startWidth, startHeight, context, context.Image().At(startWidth, startHeight), fillColor)
}

func hardFillRecurse(width int, height int, context *gg.Context, fillColor color.Color) {
	if width < 0 || height < 0 || width >= context.Width() || height >= context.Height() {
		return
	}
	currR, currG, currB, currA := context.Image().At(width, height).RGBA()
	startingR, startingG, startingB, startingA := fillColor.RGBA()
	if currR != startingR || currG != startingG || currB != startingB || currA != startingA {
		context.SetColor(fillColor)
		context.SetPixel(width, height)
		hardFillRecurse(width-1, height, context, fillColor)
		hardFillRecurse(width+1, height, context, fillColor)
		hardFillRecurse(width, height-1, context, fillColor)
		hardFillRecurse(width, height+1, context, fillColor)
	}
}

func hardFill(startWidth int, startHeight int, context *gg.Context, fillColor color.Color) {
	hardFillRecurse(startWidth, startHeight, context, fillColor)
}

func highlightBorder(context *gg.Context) {
	midwayHeight := context.Height() / 2
	for i := context.Width() - 1; i >= 0; i-- {
		r, g, b, _ := context.Image().At(i, midwayHeight).RGBA()
		if r == 0 && g == 0 && b == 0 {
			fill(i, midwayHeight, context, color.RGBA{255, 0, 0, 255})
			break
		}
	}
}

func fillEverythingOutsideBorder(context *gg.Context) {
	hardFill(0, 0, context, color.RGBA{255, 0, 0, 255})
}

func convertToBinaryColors(context *gg.Context, blackThreshold uint32) *gg.Context {
	out := gg.NewContext(context.Width(), context.Height())
	for y := 0; y < context.Height(); y++ {
		for x := 0; x < context.Width(); x++ {
			r, g, b, _ := context.Image().At(x, y).RGBA()
			if r >= blackThreshold && g >= blackThreshold && b >= blackThreshold {
				out.SetColor(color.RGBA{255, 255, 255, 255})
			} else {
				out.SetColor(color.RGBA{0, 0, 0, 255})
			}
			out.SetPixel(x, y)
		}
	}
	return out
}

// standardDeviation returns the standard deviation of the y axis
func standardDeviation(points []Point) float64 {
	n := float64(len(points))
	total := 0.0
	for i := 0; i < len(points); i++ {
		total += points[i].y
	}
	mean := total / n

	result := 0.0
	for i := 0; i < len(points); i++ {
		result += math.Pow(points[i].y-mean, 2.0)
	}
	return math.Sqrt(result / n)
}

func borderSize(context *gg.Context) (int, int) {
	left, right, top, bottom := context.Width(), 0, 0, context.Height()

	for x := 0; x < context.Width(); x++ {
		for y := 0; y < context.Height(); y++ {
			r, g, b, _ := context.Image().At(x, y).RGBA()
			if r == 65535 && g == 0 && b == 0 {
				if x < left {
					left = x
				}
				if x > right {
					right = x
				}
				if y > top {
					top = y
				}
				if y < bottom {
					bottom = y
				}
			}
		}
	}

	return right - left, top - bottom
}

func findAndHighlightPoints(context *gg.Context) []Point {
	points := make([]Point, 0)
	for y := 0; y < context.Height(); y++ {
		for x := 0; x < context.Width(); x++ {
			r, g, b, _ := context.Image().At(x, y).RGBA()
			if r == 0 && g == 0 && b == 0 {
				points = append(points, *newPoint(x, y))
				fill(x, y, context, color.RGBA{0, 255, 0, 255})
			} else if r == 65535 && ((g == 0 && b == 0) || (g == 65535 && b == 65535)) {
				context.SetColor(color.RGBA{255, 255, 255, 0})
				context.SetPixel(x, y)
			}
		}
	}
	return points
}

func add(base *gg.Context, overlap *gg.Context) *gg.Context {
	out := gg.NewContext(base.Width(), base.Height())
	for y := 0; y < out.Height(); y++ {
		for x := 0; x < out.Width(); x++ {
			_, _, _, overA := overlap.Image().At(x, y).RGBA()
			if overA > 0 {
				out.SetColor(overlap.Image().At(x, y))
			} else {
				out.SetColor(base.Image().At(x, y))
			}
			out.SetPixel(x, y)
		}
	}
	return out
}

func main() {
	var imageName = flag.String("in", "input.png", "Name of the image file to examine")
	var outImageName = flag.String("out", "out.png", "Name of the image output for human checking")
	var binaryColorThreshold = flag.Int("threshold", 24000, "Adjust this value if some of your points arn't being found")

	flag.Parse()

	if *binaryColorThreshold < 0 {
		panic("Threshold can not be below 0!")
	} else if *binaryColorThreshold > 65535 {
		panic("Threshold can not go above 65535!")
	}

	im, err := gg.LoadPNG(*imageName)
	if err != nil {
		panic("unable to load image: " + err.Error())
	}

	// The image has many gradiants of grey (nothing is truley black)
	// change to either black or white based on some threshold
	ctx := convertToBinaryColors(gg.NewContextForImage(im), uint32(*binaryColorThreshold))

	// Highlight the border to seperate what's inside it from the outside
	highlightBorder(ctx)

	// Get the height of the border for later calculations
	_, height := borderSize(ctx)

	// Clear everything that's not our graph data
	fillEverythingOutsideBorder(ctx)

	points := findAndHighlightPoints(ctx)

	log.Printf("Points found: %d", len(points))
	log.Printf("Standard Deviation: %f", standardDeviation(points)/float64(height))

	add(gg.NewContextForImage(im), ctx).SavePNG(*outImageName)
}
