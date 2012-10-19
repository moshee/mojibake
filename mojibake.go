// Attempt to decode mojibake-encoded text into readable UTF-8.
// Mojibake happens from misinterpretation and misconversion of bytes. Stuff is read in
// from Shift-JIS, for example, and the raw bytes are each interpreted as a character
// (from CP473, for example). These 1-byte values are then displayed as (converted to)
// UTF-8 multibyte characters.
package mojibake

import (
	"io"
	"strings"
	"bytes"
)

type UTF8Decoder struct {
	io.Writer
}

func NewUTF8Decoder(w io.Writer) *UTF8Decoder {
	return UTF8Decoder{w}
}

func (w *UTF8Decoder) Write(b []byte) (n int, err error) {
	return w.Writer.Write(fromUTF8(b))
}

func fromUTF8(in []byte) (out []byte) {
	var ch byte
	for i := 0; i < len(b); i++ {
		ch = b[i]
		if ch < 128 {
			out = append(out, ch)
		} else {
			i++
			out = append(out, cp473[rune(ch)<<8 | rune(b[i])])
		}
	}
}

// Use this function if the source was probably UTF-8.
func FromUTF8(s string) string {
	return string(fromUTF8([]byte(s)))
}

type SJISDecoder struct {
	io.Writer
}

func (w *SJISDecoder) Write(b []byte (n int, err error) {
	return w.Writer.Write(SJISToUTF8(fromUTF8(b)))
}

func NewSJISDecoder(w io.Writer) *SJISDecoder {
	return SJISDecoder{w}
}

// Use this function if the source was probably Shift-JIS. Panics if string is malformed.
func FromSJIS(s string) string {
	return string(SJISToUTF8(FromUTF8(s)))
}

// Use this to convert directly from Shift-JIS to UTF-8
func SJISToUTF8(s []byte) (runes []rune) {
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

	return runes
}
