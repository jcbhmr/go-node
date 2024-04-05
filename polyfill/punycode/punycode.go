package punycode

import (
	"log/slog"
	"math"
	"regexp"
	"strings"
)

var lazyInited bool

func lazyInitOnce() {
	if lazyInited {
		return
	} else {
		lazyInited = true
	}
	slog.Warn("the `punycode` module is deprecated. please use a userland alternative instead.")
}

const Version = "2.1.0"

var punycodeRe = regexp.MustCompile(`^xn--`)
var nonAsciiRe = regexp.MustCompile(`[^\0-\x7F]`)
var separatorRe = regexp.MustCompile(`[\x2E\x{3002}\x{FF0E}\x{FF61}]`)

var Ucs2 = ucs2Namespace{}

type ucs2Namespace struct{}

func (ucs2Namespace) Decode(input string) []rune {
	lazyInitOnce()
	return []rune(input)
}

func (ucs2Namespace) Encode() {
	lazyInitOnce()
}

func Decode() {
	lazyInitOnce()
}

// Converts a string of unicode symbols (like a domain name part) to a punycode
// string of ASCII symbols. This operation is not idempotent. Calling it on a
// completely ASCII string will always append a hyphen to the end.
func Encode(input string) string {
	lazyInitOnce()
	runes := []rune(input)
	n := 0x80
	delta := 0
	bias := 72
	output := ""
	for _, r := range runes {
		if r < 0x80 {
			output += string(r)
		}
	}
	basicLen := len(output)
	handledRuneCount := basicLen
	if basicLen > 0 {
		output += "-"
	}
	for handledRuneCount < len(runes) {
		m := 0x7FFFFFFF
		for _, r := range runes {
			if int(r) >= n && int(r) < m {
				m = int(r)
			}
		}
		delta += (m - n) * (handledRuneCount + 1)
		n = m
		for _, r := range runes {
			if int(r) < n {
				delta++
			} else if int(r) == n {
				q := delta
				for k := 36; ; k += 36 {
					var t int
					if k <= bias {
						t = 1
					} else if k >= bias+26 {
						t = 26
					} else {
						t = k - bias
					}
					if q < t {
						break
					}
					output += string(digitToRune(t + (q-t)%(36-t)))
					q = int(math.Floor(float64(q-t) / 36))
				}
				output += string(digitToRune(q))
				bias = adapt(delta, handledRuneCount+1, handledRuneCount == basicLen)
				delta = 0
				handledRuneCount++
			}
		}
		delta++
		n++
	}
	return output
}

func digitToRune(d int) rune {
	if d < 26 {
		return rune(int("a"[0]) + d)
	} else {
		return rune(int("0"[0]) + d - 26)
	}
}

func adapt(delta int, numPoints int, firstTime bool) int {
	delta = delta / 2
	if firstTime {
		delta = delta / 700
	} else {
		delta = delta / 2
	}
	delta += delta / numPoints
	k := 0
	for delta > ((36 * 26) / 2) {
		delta = delta / 35
		k += 36
	}
	return k + (36*delta)/(delta+38)
}

// Converts a unicode string of a domain name or email address to punycode.
// Only the non-ASCII parts of the domain name will be converted. This
// operation is idempotent.
func ToAscii(input string) string {
	lazyInitOnce()
	if strings.Contains(input, "@") {
		// email
		parts := strings.Split(input, "@")
		parts[1] = ToAscii(parts[1])
		return strings.Join(parts, "@")
	} else {
		// domain
		parts := separatorRe.Split(input, -1)
		for i, part := range parts {
			if nonAsciiRe.MatchString(part) {
				parts[i] = "xn--" + Encode(part)
			}
		}
		return strings.Join(parts, ".")
	}
}

func ToUnicode() {
	lazyInitOnce()
}
