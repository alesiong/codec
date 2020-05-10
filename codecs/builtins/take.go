package builtins

import (
	"errors"
	"io"
	"strconv"

	"github.com/alesiong/codec/codecs"
)

func init() {
	codecs.Register("take", takeCodecs{})
}

type takeCodecs struct {
}

func (s takeCodecs) RunCodec(input io.Reader, globalMode codecs.CodecMode, options map[string]string, output io.Writer) (err error) {
	takeBytes := 0

	if options["B"] != "" {
		takeBytes, err = strconv.Atoi(options["B"])
		if err != nil {
			return
		}
		if takeBytes < 0 {
			return errors.New("take: take bytes count cannot be minus")
		}
	}
	_, err = io.CopyN(output, input, int64(takeBytes))
	if err != nil && !errors.Is(err, io.EOF) {
		return
	}
	return nil
}
