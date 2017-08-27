package gomoji

import (
	"strings"
	"testing"
)

func TestFormat(t *testing.T) {
	if txt := Format(":beer:_:heart:_:computer:_:email:::"); txt != "🍺_❤️_💻_✉️::" {
		t.Fatal("err: Formatting failed: " + txt)
	}
}

func BenchmarkFormat(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Format(strings.Repeat(":computer::heart:", 1000))
	}
}
