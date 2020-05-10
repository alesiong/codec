package builtins

import (
	"errors"
	"io"
	"io/ioutil"
	"strconv"

	"github.com/alesiong/codec/codecs"
)

func init() {
	codecs.Register("drop", dropCodecs{})
}

type dropCodecs struct {
}

func (d dropCodecs) RunCodec(input io.Reader, globalMode codecs.CodecMode, options map[string]string, output io.Writer) (err error) {
	dropBytes := 0
	if options["B"] != "" {
		dropBytes, err = strconv.Atoi(options["B"])
		if err != nil {
			return
		}
		if dropBytes < 0 {
			return errors.New("drop: drop bytes count cannot be minus")
		}
	}

	_, err = io.CopyN(ioutil.Discard, input, int64(dropBytes))
	if err != nil && !errors.Is(err, io.EOF) {
		return
	}

	_, err = io.Copy(output, input)
	if err != nil {
		return
	}

	return
}
