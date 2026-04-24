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
