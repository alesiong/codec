package codecs

import (
	"bytes"
	"encoding/hex"
	"errors"
	"io"
)

type hexCodecs struct {
}

func (h hexCodecs) Execute(input io.Reader, globalMode CodecMode, options map[string]string, output io.WriteCloser) (err error) {
	useCapital := false

	if options["c"] != "" {
		useCapital = true
	}

	switch globalMode {
	case CodecModeEncoding:
		var encoder io.Writer

		if useCapital {
			encoder = hex.NewEncoder(&capitalWriter{output})
		} else {
			encoder = hex.NewEncoder(output)
		}
		err = ReadToWriter(input, encoder, output)

	case CodecModeDecoding:
		decoder := hex.NewDecoder(input)
		err = ReadToWriter(decoder, output, output)

	default:
		return errors.New("invalid codec mode")
	}

	return
}

type capitalWriter struct {
	writer io.Writer
}

func (c *capitalWriter) Write(p []byte) (n int, err error) {
	return c.writer.Write(bytes.ToUpper(p))
}
