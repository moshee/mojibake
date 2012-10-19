// Attempt to decode mojibake-encoded text into readable UTF-8.
// Mojibake happens from misinterpretation and misconversion of bytes. Stuff is read in
// from Shift-JIS, for example, and the raw bytes are each interpreted as a character
// (from CP473, for example). These 1-byte values are then displayed as (converted to)
// UTF-8 multibyte characters.
package mojibake

import "io"

type (
	UTF8Decoder struct { io.Writer }
	SJISDecoder struct { io.Writer }
)

func (d UTF8Decoder) Write(b []byte) (n int, err error) {
	return d.Writer.Write([]byte(FromUTF8(string(b))))
}

func (d SJISDecoder) Write(b []byte) (n int, err error) {
	return d.Writer.Write([]byte(FromSJIS(string(b))))
}


// Use this function if the source was probably UTF-8.
func FromUTF8(s string) string {
	bytes := make([]byte, 0)
	var ascii byte
	for _, ch := range s {
		if ascii = cp473[ch]; ascii != 0 {
			bytes = append(bytes, ascii)
		} else {
			bytes = append(bytes, byte(ch))
		}
	}
	return string(bytes)
}

// Use this function if the source was probably Shift-JIS. Panics if string is malformed.
func FromSJIS(s string) string {
	return SJISToUTF8(FromUTF8(s))
}

// Use this to convert directly from Shift-JIS to UTF-8
func SJISToUTF8(s string) string {
	runes := make([]rune, 0)
	var ch rune
	for i := 0; i < len(s); i++ {
		ch = rune(s[i])
		if ch < 128 {
			runes = append(runes, ch)
		} else {
			i++
			runes = append(runes, cp932[ch<<8 | rune(s[i])])
		}
	}

	return string(runes)
}
