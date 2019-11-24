package codecs

import (
	"compress/zlib"
	"errors"
	"io"
	"strconv"
)

type zlibCodec struct {
}

func (b zlibCodec) Execute(input io.Reader, globalMode CodecMode, options map[string]string, output io.WriteCloser) (err error) {
	level := -1
	if options["L"] != "" {
		level, err = strconv.Atoi(options["L"])
		if err != nil {
			return
		}
	}

	switch globalMode {
	case CodecModeEncoding:
		var w *zlib.Writer
		w, err = zlib.NewWriterLevel(output, level)
		if err != nil {
			return
		}
		err = ReadToWriter(input, w, w)
		if err != nil {
			return
		}
		err = output.Close()
		if err != nil {
			return
		}

		return
	case CodecModeDecoding:
		var r io.ReadCloser
		r, err = zlib.NewReader(input)
		if err != nil {
			return
		}
		err = ReadToWriter(r, output, output)
		if err != nil {
			return
		}

		err = r.Close()
		if err != nil {
			return
		}

		return
	default:
		return errors.New("invalid codec mode")
	}
}
