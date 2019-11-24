package codecs

import (
	"errors"
	"io"
)

type appendCodecs struct {
}

func (a appendCodecs) RunCodec(input io.Reader, globalMode CodecMode, options map[string]string, output io.Writer) (err error) {
	value := options["A"]
	if value == "" {
		return errors.New("append: missing required option append value (-A)")
	}

	err = ReadToWriter(input, output)
	if err != nil {
		return
	}
	_, err = output.Write([]byte(value))
	return
}

type newLineCodecs struct {
	appendCodecs
}

func (a newLineCodecs) RunCodec(input io.Reader, globalMode CodecMode, options map[string]string, output io.Writer) (err error) {
	return a.appendCodecs.RunCodec(input, globalMode, map[string]string{
		"A": "\n",
	}, output)
}
