package draw

import (
	"github.com/dihedron/overlay/command/draw/canvas"
	"github.com/dihedron/overlay/command/draw/circle"
	"github.com/dihedron/overlay/command/draw/ellipse"
	"github.com/dihedron/overlay/command/draw/image"
	"github.com/dihedron/overlay/command/draw/rectangle"
	"github.com/dihedron/overlay/command/draw/text"
)

// Commands is the set of root command groups.
type Commands struct {
	// Canvas creates a new image with the given size and colour.
	Canvas canvas.Canvas `command:"canvas" alias:"c" description:"Create a new image with the given size and colour." `
	// Rectangle adds a rectangle as an overlay to an image.
	Rectangle rectangle.Rectangle `command:"rectangle" alias:"r" description:"Add a rectangle as an overlay to an image." `
	// Circle adds a circle as an overlay to an image.
	Circle circle.Circle `command:"circle" alias:"o" description:"Add a circle as an overlay to an image." `
	// Ellipse adds an ellipse as an overlay to an image.
	Ellipse ellipse.Ellipse `command:"ellipse" alias:"e" description:"Add an ellipse as an overlay to an image." `
	// Image superimposes an image as an overlay to the given image.
	Image image.Image `command:"image" alias:"i" description:"Superimposes an image as an overlay to the given image." `
	// Text adds text as an overlay to an image.
	Text text.Text `command:"text" alias:"t" description:"Add text as an overlay to an image." `
}
