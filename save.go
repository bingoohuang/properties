package properties

import (
	"bytes"
	"fmt"
	"io"
)

// String gives the whole properties as a string
func (p Doc) String() string {
	s, _ := p.Export()

	return s
}

// Save saves the doc to file or stream.
func (p Doc) Export() (string, error) {
	buf := bytes.NewBufferString("")
	err := p.Save(buf)

	return buf.String(), err
}

// Save saves the doc to file or stream.
func (p Doc) Save(writer io.Writer) error {
	var err error

	p.Accept(func(typo byte, value string, key string) bool {
		if isProperty(typo) {
			_, err = fmt.Fprintf(writer, "%s%c%s\n", key, typo, value)
		} else {
			_, err = fmt.Fprintln(writer, value)
		}

		return nil == err
	})

	return err
}
