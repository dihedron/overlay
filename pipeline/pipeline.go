package pipeline

import (
	"image/color"

	"github.com/gogpu/gg"
)

// Canvas is the background canvas on which the image is rendered.
type Canvas struct {
	sizeX   int
	sizeY   int
	context *gg.Context
}

// NewCanvas creates a new canvas with the given size.
func NewCanvas(sizeX, sizeY int) *Canvas {
	return &Canvas{
		sizeX:   sizeX,
		sizeY:   sizeY,
		context: gg.NewContext(sizeX, sizeY),
	}
}

// Close releases the resources used by the canvas.
func (c *Canvas) Close() error {
	if c != nil && c.context != nil {
		c.context.Close()
	}
	return nil
}

type Painter func(c *Canvas) error

func (c *Canvas) Apply(painters ...Painter) (*Canvas, error) {
	for _, paint := range painters {
		if err := paint(c); err != nil {
			return c, err
		}
	}
	return c, nil
}

func Backdrop(color color.Color) Painter {
	return func(c *Canvas) error {
		r, g, b, a := color.RGBA()
		c.context.ClearWithColor(gg.RGBA{
			R: float64(r),
			G: float64(g),
			B: float64(b),
			A: float64(a),
		})
		return nil
	}
}
