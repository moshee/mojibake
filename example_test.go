package mojibake

import (
	"os"
	"io"
	"fmt"
)

func ExampleDecoder() {
	decoder, _ := NewDecoder(os.Stdout, CP473, CP932)

	// The decoding happens here
	io.Copy(decoder, os.Stdin)

	// The write to stdout happens here
	_, err := decoder.Flush()

	if err != nil {
		fmt.Println(err)
	}

	// don't forget to close it
	decoder.Close()
}
