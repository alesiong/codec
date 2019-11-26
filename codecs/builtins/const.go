package builtins

import (
	"bytes"
	"errors"
	"io"

	"github.com/alesiong/codec/codecs"
)

type constCodecs struct {
	idCodecs
}

func (c constCodecs) RunCodec(input io.Reader, globalMode codecs.CodecMode, options map[string]string, output io.Writer) (err error) {
	value := options["C"]
	if value == "" {
		return errors.New("const: missing required option const value (-C)")
	}

	valueReader := bytes.NewReader([]byte(value))

	return c.idCodecs.RunCodec(valueReader, globalMode, options, output)
}
