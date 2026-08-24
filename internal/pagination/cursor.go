package pagination

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
)

type Cursor struct {
	Offset int
	Limit  int
}

func Encode(c Cursor) string {
	return base64.RawURLEncoding.EncodeToString([]byte(fmt.Sprintf("%d:%d", c.Offset, c.Limit)))
}
func Decode(v string) (Cursor, error) {
	if v == "" {
		return Cursor{}, nil
	}
	b, err := base64.RawURLEncoding.DecodeString(v)
	if err != nil {
		return Cursor{}, err
	}
	parts := strings.Split(string(b), ":")
	if len(parts) != 2 {
		return Cursor{}, fmt.Errorf("invalid cursor")
	}
	o, err := strconv.Atoi(parts[0])
	if err != nil {
		return Cursor{}, err
	}
	l, err := strconv.Atoi(parts[1])
	if err != nil {
		return Cursor{}, err
	}
	if o < 0 || l <= 0 || l > 100 {
		return Cursor{}, fmt.Errorf("cursor bounds")
	}
	return Cursor{o, l}, nil
}
func Next(c Cursor, count int) string {
	if count < c.Limit {
		return ""
	}
	return Encode(Cursor{Offset: c.Offset + c.Limit, Limit: c.Limit})
}
