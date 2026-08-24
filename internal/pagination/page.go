package pagination

import "strconv"

type Page struct {
	Limit  int
	Offset int
}

func Parse(limit, offset string) Page {
	l, _ := strconv.Atoi(limit)
	o, _ := strconv.Atoi(offset)
	if l <= 0 || l > 100 {
		l = 25
	}
	if o < 0 {
		o = 0
	}
	return Page{l, o}
}
