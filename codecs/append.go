package codecs

import (
	"errors"
	"io"
)

type appendCodecs struct {
}

func (a appendCodecs) Execute(input io.Reader, globalMode CodecMode, options map[string]string, output io.WriteCloser) (err error) {
	value := options["A"]
	if value == "" {
		return errors.New("append: missing required option append value (-A)")
	}

	err = ReadToWriter(input, output, nil)
	if err != nil {
		return
	}
	_, err = output.Write([]byte(value))
	if err != nil {
		return
	}
	err = output.Close()
	return
}

type newLineCodecs struct {
	appendCodecs
}

func (a newLineCodecs) Execute(input io.Reader, globalMode CodecMode, options map[string]string, output io.WriteCloser) (err error) {
	return a.appendCodecs.Execute(input, globalMode, map[string]string{
		"A": "\n",
	}, output)
}
