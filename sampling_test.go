package jigar

import "testing"

func TestDecideShouldTrace(t *testing.T) {
	cases := []struct {
		name   string
		header string
		ratio  float64
		want   bool
	}{
		{"header forces on", "true", 0, true},
		{"header empty, ratio 0 -> off", "", 0, false},
		{"header empty, ratio 1 -> on", "", 1, true},
		{"header false, ratio 1 -> on (header only forces 'true')", "false", 1, true},
		{"header empty, negative ratio -> off", "", -0.5, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := DecideShouldTrace(c.header, c.ratio)
			if got != c.want {
				t.Fatalf("DecideShouldTrace(%q, %v) = %v, want %v", c.header, c.ratio, got, c.want)
			}
		})
	}
}
