package builtins

import (
	"bytes"
	"errors"
	"io"
	"strconv"

	"github.com/alesiong/codec/codecs"
)

func init() {
	codecs.Register("repeat", repeatCodecs{})
	codecs.Register("id", idCodecs{})
}

type repeatCodecs struct {
}

func (r repeatCodecs) Usage() string {
	return "    -T times: repeat input for `times` times (int, >=0, default 0)"
}

func (r repeatCodecs) RunCodec(input io.Reader, globalMode codecs.CodecMode, options map[string]string, output io.Writer) (err error) {
	times := 0
	if options["T"] != "" {
		times, err = strconv.Atoi(options["T"])
		if err != nil {
			return
		}
		if times < 0 {
			return errors.New("repeat: times cannot be minus")
		}
	}

	var buffer bytes.Buffer

	writer := io.MultiWriter(output, &buffer)

	for i := 0; i < times; i++ {
		if i == 0 {
			_, err = io.Copy(writer, input)
		} else {
			_, err = io.Copy(output, bytes.NewReader(buffer.Bytes()))
		}
		if err != nil {
			return
		}
	}
	return
}

type idCodecs struct {
	repeatCodecs
}

func (i idCodecs) Usage() string {
	return "    pass input to output as is"
}

func (i idCodecs) RunCodec(input io.Reader, globalMode codecs.CodecMode, options map[string]string, output io.Writer) (err error) {
	return i.repeatCodecs.RunCodec(input, globalMode, map[string]string{"T": "1"}, output)
}
