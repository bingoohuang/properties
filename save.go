package properties

import (
	"bytes"
	"fmt"
	"io"
)

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
		switch typo {
		case '#', '!', ' ':
			_, err = fmt.Fprintln(writer, value)
		case '=', ':':
			_, err = fmt.Fprintf(writer, "%s%c%s\n", key, typo, value)
		}

		return nil == err
	})

	return err
}
