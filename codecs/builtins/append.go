package builtins

import (
	"errors"
	"io"

	"github.com/alesiong/codec/codecs"
)

type appendCodecs struct {
}

func (a appendCodecs) RunCodec(input io.Reader, globalMode codecs.CodecMode, options map[string]string, output io.Writer) (err error) {
	value := options["A"]
	if value == "" {
		return errors.New("append: missing required option append value (-A)")
	}

	_, err = io.Copy(output, input)
	if err != nil {
		return
	}
	_, err = output.Write([]byte(value))
	return
}

type newLineCodecs struct {
	appendCodecs
}

func (a newLineCodecs) RunCodec(input io.Reader, globalMode codecs.CodecMode, options map[string]string, output io.Writer) (err error) {
	return a.appendCodecs.RunCodec(input, globalMode, map[string]string{
		"A": "\n",
	}, output)
}
