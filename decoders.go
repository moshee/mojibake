package mojibake

import (
	"errors"
	"unicode/utf8"
)

// Assume CP473 ASCII source mangled to UTF-8. Multibyte to single byte
func dec_cp473(in, out chan byte, finished, closed chan error) {
	var (
		b    byte
		ok   bool
		ch   rune
		full bool
		dec  byte
		cp   = make([]byte, 0, utf8.UTFMax)
	)

loop:
	for {
		select {
		case b, ok = <-in:
			if !ok {
				break loop
			}
			cp = append(cp, b)

			if full = utf8.FullRune(cp); !full {
				continue loop
			}

			ch, _ = utf8.DecodeRune(cp)
			// XXX: change this to cp[0:0:utf8.UTFMax] for go1.2 slice syntax
			cp = cp[0:0]

			if dec = cp473[ch]; dec != 0 {
				out <- dec
			} else {
				out <- byte(ch)
			}

		case <-finished:
			if !full {
				finished <- errors.New("mojibake: dec_cp473: malformed byte stream")
			} else {
				finished <- nil
			}

		case <-closed:
			break loop
		}
	}
}

// Assume Shift-JIS source mangled to UTF-8
func dec_cp932(in, out chan byte, finished, closed chan error) {
	dec_table(in, out, finished, closed, cp932[:])
}

// Assume CJK source mangled to UTF-8
func dec_cp936(in, out chan byte, finished, closed chan error) {
	dec_table(in, out, finished, closed, cp936[:])
}

// Generalized call for multibyte (2 max) decoding (SJIS and CJK). Multibyte to multibyte
func dec_table(in, out chan byte, finished, closed chan error, table []rune) {
	var (
		b     byte
		c     byte
		ok    bool
		multi bool
		sz    int
		cp    = make([]byte, utf8.UTFMax)
	)

loop:
	for {
		select {
		case b, ok = <-in:
			if !ok {
				break loop
			}

			multi = false

			if b < 128 {
				utf8.EncodeRune(cp, rune(b))
				out <- cp[0]
			} else {
				multi = true
				if c, ok = <-in; !ok {
					break loop
				}

				multi = false

				sz = utf8.EncodeRune(cp, table[rune(b)<<8|rune(c)])
				// we don't really care about garbage in the slice
				for _, encoded := range cp[:sz] {
					out <- encoded
				}
			}

		case <-finished:
			if multi {
				finished <- errors.New("mojibake: dec_table: malformed byte stream")
			} else {
				finished <- nil
			}

		case <-closed:
			break loop
		}
	}
}
