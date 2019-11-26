package builtins

import (
	"encoding/base64"
	"errors"
	"io"

	"github.com/alesiong/codec/codecs"
)

type base64Codec struct {
}

func (b base64Codec) RunCodec(input io.Reader, globalMode codecs.CodecMode, options map[string]string, output io.Writer) (err error) {
	encoding := base64.StdEncoding
	if options["u"] != "" {
		encoding = base64.URLEncoding
	}

	switch globalMode {
	case codecs.CodecModeEncoding:
		encoder := base64.NewEncoder(encoding, output)
		err = codecs.ReadToWriter(input, encoder)
		if err != nil {
			return
		}
		err = encoder.Close()
		return
	case codecs.CodecModeDecoding:
		decoder := base64.NewDecoder(encoding, input)
		err = codecs.ReadToWriter(decoder, output)
		return
	default:
		return errors.New("invalid codec mode")
	}
}
