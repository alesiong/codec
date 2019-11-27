package builtins

import (
	"bytes"
	"encoding/hex"
	"errors"
	"io"

	"github.com/alesiong/codec/codecs"
)

type hexCodecs struct {
}

func (h hexCodecs) RunCodec(input io.Reader, globalMode codecs.CodecMode, options map[string]string, output io.Writer) (err error) {
	useCapital := false

	if options["c"] != "" {
		useCapital = true
	}

	switch globalMode {
	case codecs.CodecModeEncoding:
		var encoder io.Writer

		if useCapital {
			encoder = hex.NewEncoder(&capitalWriter{output})
		} else {
			encoder = hex.NewEncoder(output)
		}
		_, err = io.Copy(encoder, input)

	case codecs.CodecModeDecoding:
		decoder := hex.NewDecoder(input)
		_, err = io.Copy(output, decoder)

	default:
		return errors.New("invalid codec mode")
	}

	return
}

type capitalWriter struct {
	writer io.Writer
}

func (c *capitalWriter) Write(p []byte) (n int, err error) {
	n = len(p)
	_, err = c.writer.Write(bytes.ToUpper(p))
	return
}
