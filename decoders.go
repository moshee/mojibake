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
			cp = cp[0:0]
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
	dec_table(in, out, finished, closed, CP932)
}

// Assume CJK source mangled to UTF-8
func dec_cp936(in, out chan byte, finished, closed chan error) {
	dec_table(in, out, finished, closed, CP936)
}

// Generalized call for multibyte (2 max) decoding (SJIS and CJK). Multibyte to multibyte
func dec_table(in, out chan byte, finished, closed chan error, enc Encoding) {
	table := enc_tables[enc]

	var (
		b     byte
		a     byte
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

			if multi {
				sz = utf8.EncodeRune(cp, table[rune(a)<<8|rune(b)])
				// we don't really care about garbage in the slice
				for _, encoded := range cp[:sz] {
					out <- encoded
				}
			} else {
				if b < 128 {
					utf8.EncodeRune(cp, rune(b))
					out <- cp[0]
				} else {
					multi = true
					a = b
					continue loop
				}
			}
			multi = false

		case <-finished:
			if multi {
				finished <- errors.New("mojibake: dec_table(" + enc.String() + "): malformed byte stream")
			} else {
				finished <- nil
			}
			multi = false

		case <-closed:
			break loop
		}
	}
}
