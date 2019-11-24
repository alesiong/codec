package codecs

import (
	"encoding/base64"
	"errors"
	"io"
)

type base64Codec struct {
}

func (b base64Codec) Execute(input io.Reader, globalMode CodecMode, options map[string]string, output io.WriteCloser) (err error) {
	encoding := base64.StdEncoding
	if options["u"] != "" {
		encoding = base64.URLEncoding
	}

	switch globalMode {
	case CodecModeEncoding:
		encoder := base64.NewEncoder(encoding, output)
		err = ReadToWriter(input, encoder, encoder)
		if err != nil {
			return
		}
		err = output.Close()
		return
	case CodecModeDecoding:
		decoder := base64.NewDecoder(encoding, input)
		err = ReadToWriter(decoder, output, output)
		return
	default:
		return errors.New("invalid codec mode")
	}
}
