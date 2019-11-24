package codecs

import (
	"bytes"
	"errors"
	"io"
)

type constCodecs struct {
	idCodecs
}

func (c constCodecs) Execute(input io.Reader, globalMode CodecMode, options map[string]string, output io.WriteCloser) (err error) {
	value := options["C"]
	if value == "" {
		return errors.New("const: missing required option const value (-C)")
	}

	valueReader := bytes.NewReader([]byte(value))

	return c.idCodecs.Execute(valueReader, globalMode, options, output)
}
