package pagination

import "testing"

func TestCursorRoundTrip(t *testing.T) {
	c := Cursor{Offset: 20, Limit: 25}
	v := Encode(c)
	got, err := Decode(v)
	if err != nil || got != c {
		t.Fatal(v, got, err)
	}
	if _, err := Decode("bad"); err == nil {
		t.Fatal("bad cursor")
	}
	if Next(c, 10) != "" {
		t.Fatal("terminal next")
	}
	if Next(c, 25) == "" {
		t.Fatal("next missing")
	}
}
func TestParse(t *testing.T) {
	if p := Parse("10", "-1"); p.Limit != 10 || p.Offset != 0 {
		t.Fatal(p)
	}
	if p := Parse("0", "0"); p.Limit != 25 {
		t.Fatal(p)
	}
}
