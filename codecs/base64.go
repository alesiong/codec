package codecs

import (
	"encoding/base64"
	"errors"
	"io"
)

type base64Codec struct {
}

func (b base64Codec) RunCodec(input io.Reader, globalMode CodecMode, options map[string]string, output io.Writer) (err error) {
	encoding := base64.StdEncoding
	if options["u"] != "" {
		encoding = base64.URLEncoding
	}

	switch globalMode {
	case CodecModeEncoding:
		encoder := base64.NewEncoder(encoding, output)
		err = ReadToWriter(input, encoder)
		if err != nil {
			return
		}
		err = encoder.Close()
		return
	case CodecModeDecoding:
		decoder := base64.NewDecoder(encoding, input)
		err = ReadToWriter(decoder, output)
		return
	default:
		return errors.New("invalid codec mode")
	}
}
