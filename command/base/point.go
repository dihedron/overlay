package base

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Size is a 2D coordinate.
type Size struct {
	X, Y int
}

// UnmarshalFlag parses a string representation of a size in the format "x,y".
func (p *Size) UnmarshalFlag(value string) error {
	parts := strings.Split(value, ",")

	if len(parts) != 2 {
		return errors.New("invalid format: expected two numbers separated by a ,")
	}
	x, err := strconv.ParseInt(parts[0], 10, 32)
	if err != nil {
		return err
	}
	y, err := strconv.ParseInt(parts[1], 10, 32)
	if err != nil {
		return err
	}
	p.X = int(x)
	p.Y = int(y)
	return nil
}

// MarshalFlag returns the string representation of a size in the format "x,y".
func (p Size) MarshalFlag() (string, error) {
	return fmt.Sprintf("%d,%d", p.X, p.Y), nil
}

// Rectangle is a 2D rectangle; TopLeft and BottomRight are two opposite corners.
type Rectangle struct {
	TopLeft     Size
	BottomRight Size
}

// UnmarshalFlag parses a string representation of a rectangle in the format "x0,y0,x1,y1"
// It uses flag.Value interface so that it can be used as a flag in the command line.
func (r *Rectangle) UnmarshalFlag(value string) error {
	parts := strings.Split(value, ",")

	if len(parts) != 4 {
		return errors.New("invalid format: expected four numbers separated by a ,")
	}
	x0, err := strconv.ParseInt(parts[0], 10, 32)
	if err != nil {
		return err
	}
	y0, err := strconv.ParseInt(parts[1], 10, 32)
	if err != nil {
		return err
	}
	x1, err := strconv.ParseInt(parts[2], 10, 32)
	if err != nil {
		return err
	}
	y1, err := strconv.ParseInt(parts[3], 10, 32)
	if err != nil {
		return err
	}
	r.TopLeft = Size{X: int(x0), Y: int(y0)}
	r.BottomRight = Size{X: int(x1), Y: int(y1)}
	return nil
}

// MarshalFlag returns the string representation of a rectangle in the format "x0,y0,x1,y1".
func (r Rectangle) MarshalFlag() (string, error) {
	return fmt.Sprintf("%d,%d,%d,%d", r.TopLeft.X, r.TopLeft.Y, r.BottomRight.X, r.BottomRight.Y), nil
}

// Point is a 2D coordinate as floats.
type Point struct {
	X, Y float64
}

// UnmarshalFlag parses a string representation of a point in the format "x,y".
func (p *Point) UnmarshalFlag(value string) error {
	parts := strings.Split(value, ",")

	if len(parts) != 2 {
		return errors.New("invalid format: expected two numbers separated by a ,")
	}
	x, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return err
	}
	y, err := strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return err
	}
	p.X = x
	p.Y = y
	return nil
}

// MarshalFlag returns the string representation of a point in the format "x,y".
func (p Point) MarshalFlag() (string, error) {
	return fmt.Sprintf("%g,%g", p.X, p.Y), nil
}
