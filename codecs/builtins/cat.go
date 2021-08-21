package builtins

import (
	"errors"
	"io"
	"os"

	"github.com/alesiong/codec/codecs"
)

type catCodecs struct {
}

func init() {
	codecs.Register("cat", catCodecs{})
}

func (c catCodecs) RunCodec(input io.Reader, globalMode codecs.CodecMode, options map[string]string, output io.Writer) (err error) {
	inputFile := options["F"]
	if inputFile == "" {
		return errors.New("cat: missing required option input file (-F)")
	}

	if options["c"] == "" {
		_, err = io.Copy(output, input)
		if err != nil {
			return
		}
	}
	file, err := os.Open(inputFile)
	if err != nil {
		return
	}
	defer file.Close()
	_, err = io.Copy(output, file)
	return
}
