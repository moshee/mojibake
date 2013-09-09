// Package mojibake provides facilities to attempt to decode character encoding
// mess-ups caused by text passing through multiple encoding environments, the
// bytes being misinterpreted in each. http://en.wikipedia.org/wiki/Mojibake
package mojibake

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"unicode/utf8"
)

type Encoding int

const (
	CP473 Encoding = iota // assume UTF-8 source
	CP932                 // assume Shift-JIS source
	CP936                 // assume CJK source
)

// in:       bytes into the decode worker.
// out:      decoded bytes out of the decoder.
// finished: recieve indicates that flush has been called. Once this happens,
//           send back the error encountered or nil if none.
// closed:   recieve indicates that close has been called.
type decode_func func(in, out chan byte, finished, closed chan error)

var funcs = []decode_func{dec_cp473, dec_cp932, dec_cp936}

// Error type to signal decode_workers that they should break their loops and exit
type finished struct{}

func (self finished) Error() string {
	return "Finished encoding"
}

type decoder interface {
	chans() chan byte
}

// Write satisfies the io.Writer interface. However, due to the variability of
// character encoding widths, it will not actually write to the underlying
// io.Writer until a call to Flush—Decoder doesn't work like a regular filter.
// Instead, it keeps an internal buffer (in case e.g. buffered I/O is being
// used to write into the Decoder and the buffer boundary falls in the middle
// of a multibyte character) which is written to the underlying io.Writer
// once Flush is called.
//
// Close will call Flush if it has not been called beforehand. It is an error
// to call Close on a closed Decoder. Flush will not call Close.
type Decoder interface {
	io.WriteCloser

	// Flush triggers a Write to the underlying io.Writer.
	Flush() (written int64, err error)
}

// the root decoder. Reads bytes and runes from Write input.
type decoder_chain struct {
	w       io.Writer
	workers []*decode_worker
	buf     *bytes.Buffer
	into    chan byte
	closed  chan error
}

// Initialize a Decoder that will use the specified Encoding path to decode
// writes destined for w.
func NewDecoder(w io.Writer, encs ...Encoding) (Decoder, error) {
	if len(encs) == 0 {
		return nil, errors.New("mojibake: no encoding path specified")
	}

	chain := &decoder_chain{
		w:       w,
		workers: make([]*decode_worker, len(encs)),
		buf:     new(bytes.Buffer),
		into:    make(chan byte),
		closed:  make(chan error),
	}

	var d decoder = chain

	for i, enc := range encs {
		worker := &decode_worker{
			upstream: d,
			decode:   funcs[enc],
			out:      make(chan byte, utf8.UTFMax),
			finished: make(chan error),
		}

		chain.workers[i] = worker
		d = worker

		go worker.start(chain.closed)
	}

	go chain.start()

	return chain, nil
}

func (self *decoder_chain) Write(p []byte) (int, error) {
	for _, b := range p {
		self.into <- b
	}

	return len(p), nil
}

func (self *decoder_chain) Flush() (int64, error) {
	// send finished signal to each worker sequentially
	for _, worker := range self.workers {
		worker.finished <- finished{}
		// waiting for the recv back on the same channel ensures they don't
		// overlap
		if err := <-worker.finished; err != nil {
			self.buf.Reset()
			return 0, err
		}
	}

	return io.Copy(self.w, self.buf)
}

func (self *decoder_chain) Close() error {
	var err error
	if self.buf.Len() > 0 {
		_, err = self.Flush()
	}

	for _ = range self.workers {
		// close each worker
		self.closed <- nil
	}

	// close self
	self.closed <- nil

	return err
}

func (self *decoder_chain) start() {
	var (
		last_worker = self.workers[len(self.workers)-1]
		b           byte
		output      = last_worker.chans()
	)

	for {
		select {
		case b = <-output:
			self.buf.WriteByte(b)

		case <-self.closed:
			return
		}
	}

}

func (self *decoder_chain) chans() chan byte {
	return self.into
}

type decode_worker struct {
	upstream decoder
	decode   decode_func
	out      chan byte
	finished chan error
}

// Start the decoding loop
func (self *decode_worker) start(closed chan error) {
	in := self.upstream.chans()
	self.decode(in, self.out, self.finished, closed)
}

func (self *decode_worker) chans() chan byte {
	return self.out
}

// Like Decode, but panics if an error is found.
func MustDecode(garbled string, encs ...Encoding) string {
	dec, err := Decode(garbled, encs...)
	if err != nil {
		panic(err)
	}

	return dec
}

// TODO: heuristics

// Attempt to decode a garbled string using a specified encoding path. The
// encoding path is your guess as to what encoding environments the string
// probably passed through before reaching its current state. Note that not all
// encoding-garbled strings have recoverable data.
func Decode(garbled string, encs ...Encoding) (string, error) {
	var (
		buf      = new(bytes.Buffer)
		dec, err = NewDecoder(buf, encs...)
		r        = strings.NewReader(garbled)
	)

	if err != nil {
		return "", err
	}

	io.Copy(dec, r)

	err = dec.Close()

	return buf.String(), err
}
