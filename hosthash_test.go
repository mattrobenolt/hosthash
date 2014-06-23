package hosthash

import "testing"
import "bytes"
import "encoding/hex"

func TestRules(t *testing.T) {
	var tests = []string{
		"www.example.com",
		"*.example.com",
		"www.example.*",
		`^[w]{3}\.example\.com$`,
	}

	var h *Hash
	const host = "www.example.com"
	for i, rule := range tests {
		h = New()
		h.Add(rule, rule)
		if val, ok := h.Get(host); val != rule || !ok {
			t.Errorf("%d.Get(%v) = (%v, %v), want (%v, %v)", i, host, val, ok, rule, true)
		}
	}
}

func BenchmarkExactHit(b *testing.B) {
	h := New()
	h.Add("www.example.com", 1)
	for i := 0; i < b.N; i++ {
		h.Get("www.example.com")
	}
}

func BenchmarkExactMiss(b *testing.B) {
	h := New()
	h.Add("www.example.com", 1)
	for i := 0; i < b.N; i++ {
		h.Get("derp.example.com")
	}
}

func BenchmarkPrefixHit(b *testing.B) {
	h := New()
	h.Add("*.example.com", 1)
	for i := 0; i < b.N; i++ {
		h.Get("www.example.com")
	}
}

func BenchmarkPrefixMiss(b *testing.B) {
	h := New()
	h.Add("*.example.com", 1)
	for i := 0; i < b.N; i++ {
		h.Get("www.example.org")
	}
}

func BenchmarkSuffixHit(b *testing.B) {
	h := New()
	h.Add("www.example.*", 1)
	for i := 0; i < b.N; i++ {
		h.Get("www.example.org")
	}
}

func BenchmarkSuffixMiss(b *testing.B) {
	h := New()
	h.Add("www.example.*", 1)
	for i := 0; i < b.N; i++ {
		h.Get("www.foo.com")
	}
}

func BenchmarkDefault(b *testing.B) {
	h := New()
	h.Add("_", 1)
	for i := 0; i < b.N; i++ {
		h.Get("www.example.com")
	}
}

func BenchmarkRegexHit(b *testing.B) {
	h := New()
	h.Add(`^foo\.[a-z]+\.example\.com$`, 1)
	for i := 0; i < b.N; i++ {
		h.Get("foo.www.example.com")
	}
}

func BenchmarkRegexMiss(b *testing.B) {
	h := New()
	h.Add(`^foo\.[a-z]+\.example\.com$`, 1)
	for i := 0; i < b.N; i++ {
		h.Get("foo.com")
	}
}

func BenchmarkFallToDefault(b *testing.B) {
	h := New()
	h.Add("example.com", 1)
	h.Add("www.example.com", 1)
	h.Add("*.example.com", 1)
	h.Add("_", 1)
	for i := 0; i < b.N; i++ {
		h.Get("mattrobenolt.com")
	}
}
