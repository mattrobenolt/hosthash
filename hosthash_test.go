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

func BenchmarkString(b *testing.B) {
	m := map[string]int{
		"2c26b46b68ffc68ff99b453c1d30413413422d706483bfa0f98a5e886266e7ae": 1,
	}
	for i := 0; i < b.N; i++ {
		_ = m[hex.EncodeToString([]byte{44, 38, 180, 107, 104, 255, 198, 143, 249, 155, 69, 60, 29, 48, 65, 52, 19, 66, 45, 112, 100, 131, 191, 160, 249, 138, 94, 136, 98, 102, 231, 174})]
	}
}

func BenchmarkArray(b *testing.B) {
	m := map[[32]byte]int{
		[32]byte{44, 38, 180, 107, 104, 255, 198, 143, 249, 155, 69, 60, 29, 48, 65, 52, 19, 66, 45, 112, 100, 131, 191, 160, 249, 138, 94, 136, 98, 102, 231, 174}: 1,
	}
	for i := 0; i < b.N; i++ {
		_ = m[[32]byte{44, 38, 180, 107, 104, 255, 198, 143, 249, 155, 69, 60, 29, 48, 65, 52, 19, 66, 45, 112, 100, 131, 191, 160, 249, 138, 94, 136, 98, 102, 231, 174}]
	}
}

func BenchmarkBytes(b *testing.B) {
	foo := []byte("abcdefg")
	bar := []byte("abcdefg")
	for i := 0; i < b.N; i++ {
		_ = bytes.Equal(foo, bar)
	}
}
