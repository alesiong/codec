package builtins

import (
	"errors"
	"io"
	"os"

	"github.com/alesiong/codec/codecs"
)

type catCodecs struct {
}

func (c catCodecs) RunCodec(input io.Reader, globalMode codecs.CodecMode, options map[string]string, output io.Writer) (err error) {
	inputFile := options["F"]
	if inputFile == "" {
		return errors.New("cat: missing required option output file (-F)")
	}

	if options["c"] == "" {
		err = codecs.ReadToWriter(input, output)
		if err != nil {
			return
		}
	}
	file, err := os.Open(inputFile)
	if err != nil {
		return
	}
	defer file.Close()
	return codecs.ReadToWriter(file, output)
}
