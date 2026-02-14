package base

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Point is a 2D coordinate.
type Point struct {
	X, Y int
}

// UnmarshalFlag parses a string representation of a point in the format "x,y".
func (p *Point) UnmarshalFlag(value string) error {
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

// MarshalFlag returns the string representation of a point in the format "x,y".
func (p Point) MarshalFlag() (string, error) {
	return fmt.Sprintf("%d,%d", p.X, p.Y), nil
}
