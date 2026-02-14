package base

import (
	"fmt"
	"image/color"
	"log/slog"
	"strconv"
)

// Colour is a colour in the format #RGB, #RGBA, #RRGGBB, or #RRGGBBAA.
type Colour color.RGBA

// UnmarshalFlag parses a string representation of a colour in the format #RGB, #RGBA, #RRGGBB, or #RRGGBBAA.
func (c *Colour) UnmarshalFlag(value string) error {
	if len(value) > 0 && value[0] != '#' {
		return fmt.Errorf("invalid color string format")
	}

	switch len(value) {
	case 4: // #RGB

		if r, err := strconv.ParseUint(value[1:2], 16, 8); err != nil {
			return err
		} else {
			c.R = uint8(r)
		}

		if g, err := strconv.ParseUint(value[2:3], 16, 8); err != nil {
			return err
		} else {
			c.G = uint8(g)
		}

		if b, err := strconv.ParseUint(value[3:4], 16, 8); err != nil {
			return err
		} else {
			c.B = uint8(b)
		}

		c.A = 255

	case 5: // #RGBA

		if r, err := strconv.ParseUint(value[1:2], 16, 8); err != nil {
			return err
		} else {
			c.R = uint8(r)
		}

		if g, err := strconv.ParseUint(value[2:3], 16, 8); err != nil {
			return err
		} else {
			c.G = uint8(g)
		}

		if b, err := strconv.ParseUint(value[3:4], 16, 8); err != nil {
			return err
		} else {
			c.B = uint8(b)
		}

		if a, err := strconv.ParseUint(value[4:5], 16, 8); err != nil {
			return err
		} else {
			c.A = uint8(a)
		}

	case 7: // #RRGGBB

		if r, err := strconv.ParseUint(value[1:3], 16, 8); err != nil {
			return err
		} else {
			c.R = uint8(r)
		}

		if g, err := strconv.ParseUint(value[3:5], 16, 8); err != nil {
			return err
		} else {
			c.G = uint8(g)
		}

		if b, err := strconv.ParseUint(value[5:7], 16, 8); err != nil {
			return err
		} else {
			c.B = uint8(b)
		}

		c.A = 255

	case 9: // #RRGGBBAA

		if r, err := strconv.ParseUint(value[1:3], 16, 8); err != nil {
			return err
		} else {
			c.R = uint8(r)
		}

		if g, err := strconv.ParseUint(value[3:5], 16, 8); err != nil {
			return err
		} else {
			c.G = uint8(g)
		}

		if b, err := strconv.ParseUint(value[5:7], 16, 8); err != nil {
			return err
		} else {
			c.B = uint8(b)
		}

		if a, err := strconv.ParseUint(value[7:9], 16, 8); err != nil {
			return err
		} else {
			c.A = uint8(a)
		}
	default:
		return fmt.Errorf("invalid color string format")
	}
	slog.Debug("parsed color", "red", c.R, "green", c.G, "blue", c.B, "alpha", c.A)
	return nil
}

// MarshalFlag returns the string representation of a colour in the format #RRGGBBAA.
func (c Colour) MarshalFlag() (string, error) {
	return fmt.Sprintf("#%02X%02X%02X%02X", c.R, c.G, c.B, c.A), nil
}
