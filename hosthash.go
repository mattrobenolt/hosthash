package hosthash

import (
	"errors"
	"regexp"
	"strings"
)

var (
	ERR_TOO_SHORT = errors.New("hash: must be more than 0 characters")
	ERR_INVALID   = errors.New("hash: invalid key name")
	ERR_DUPLICATE = errors.New("hash: duplicate key")
)

type pattern struct {
	p     *regexp.Regexp
	value interface{}
}

type Hash struct {
	exact    map[string]interface{}
	prefix   map[string]interface{}
	suffix   map[string]interface{}
	patterns []*pattern
	default_ interface{}
}

func New() *Hash {
	return &Hash{}
}

func (h *Hash) Add(key string, value interface{}) error {
	l := len(key)

	if l == 0 {
		return ERR_TOO_SHORT
	}

	// Check for a single invalid character
	if l == 1 {
		switch key[0] {
		case '.', '*':
			return ERR_INVALID
		case '_':
			if h.default_ != nil {
				return ERR_DUPLICATE
			}
			h.default_ = value
			return nil
		}
	} else {

		// Regular expression
		if key[0] == '^' {
			p := regexp.MustCompile(key)
			h.patterns = append(h.patterns, &pattern{p, value})
			return nil
		}

		// Check for consecutive dots
		if strings.Index(key, "..") > -1 {
			return ERR_INVALID
		}

		if l > 2 {
			// Check for a * that isn't at the head or tail
			if strings.IndexRune(key[1:l-1], '*') > -1 {
				return ERR_INVALID
			}

			if strings.HasPrefix(key, "*.") {
				if h.prefix == nil {
					h.prefix = make(map[string]interface{})
				}
				key = key[2:]
				if _, ok := h.prefix[key]; ok {
					return ERR_DUPLICATE
				}
				h.prefix[key] = value
				return nil
			}

			if strings.HasSuffix(key, ".*") {
				if h.suffix == nil {
					h.suffix = make(map[string]interface{})
				}
				key = key[:l-2]
				if _, ok := h.suffix[key]; ok {
					return ERR_DUPLICATE
				}
				h.suffix[key] = value
				return nil
			}
		}
	}

	if h.exact == nil {
		h.exact = make(map[string]interface{})
	}

	// Must be an exact match
	if _, ok := h.exact[key]; ok {
		return ERR_DUPLICATE
	}
	h.exact[key] = value
	return nil
}

func (h *Hash) Get(key string) (value interface{}, ok bool) {
	// Check first for an exact match
	if h.exact != nil {
		if value, ok = h.exact[key]; ok {
			return
		}
	}

	l := len(key)

	if h.prefix != nil {
		if value, ok = h.getPrefix(key, l-1); ok {
			return
		}
	}

	if h.suffix != nil {
		if value, ok = h.getSuffix(key, 0, l); ok {
			return
		}
	}

	if h.patterns != nil {
		if value, ok = h.getPattern(key); ok {
			return
		}
	}

	if h.default_ != nil {
		return h.default_, true
	}

	return
}

// Match *.example.com
func (h *Hash) getPrefix(key string, start int) (value interface{}, ok bool) {
	for ; start > 0; start-- {
		if key[start] == '.' {
			break
		}
	}

	// Exhausted all possibilities
	if start == 0 {
		return nil, false
	}

	if value, ok = h.prefix[key[start+1:]]; ok {
		return
	}
	return h.getPrefix(key, start-1)
}

// Match example.*
func (h *Hash) getSuffix(key string, start, end int) (value interface{}, ok bool) {
	for ; start < end; start++ {
		if key[start] == '.' {
			break
		}
	}

	// Exhausted all possibilities
	if start == end {
		return nil, false
	}

	if value, ok = h.suffix[key[:start]]; ok {
		return
	}
	return h.getSuffix(key, start+1, end)
}

// Match against the list of patterns
// This will always be an O(n) where n is the number of patterns
func (h *Hash) getPattern(key string) (value interface{}, ok bool) {
	for _, p := range h.patterns {
		if p.p.MatchString(key) {
			return p.value, true
		}
	}
	return nil, false
}
