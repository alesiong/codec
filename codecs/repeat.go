package codecs

import (
	"bytes"
	"errors"
	"io"
	"strconv"
)

type repeatCodecs struct {
}

func (r repeatCodecs) Execute(input io.Reader, globalMode CodecMode, options map[string]string, output io.WriteCloser) (err error) {
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
			err = ReadToWriter(input, writer, nil)
		} else {
			err = ReadToWriter(bytes.NewReader(buffer.Bytes()), output, nil)
		}
	}
	err = output.Close()
	return
}

type idCodecs struct {
	repeatCodecs
}

func (i idCodecs) Execute(input io.Reader, globalMode CodecMode, options map[string]string, output io.WriteCloser) (err error) {
	return i.repeatCodecs.Execute(input, globalMode, map[string]string{"T": "1"}, output)
}
