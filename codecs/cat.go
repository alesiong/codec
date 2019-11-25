package codecs

import (
	"errors"
	"io"
	"os"
)

type catCodecs struct {
}

func (c catCodecs) RunCodec(input io.Reader, globalMode CodecMode, options map[string]string, output io.Writer) (err error) {
	inputFile := options["F"]
	if inputFile == "" {
		return errors.New("cat: missing required option output file (-F)")
	}

	if options["c"] == "" {
		err = ReadToWriter(input, output)
		if err != nil {
			return
		}
	}
	file, err := os.Open(inputFile)
	if err != nil {
		return
	}
	defer file.Close()
	return ReadToWriter(file, output)
}
