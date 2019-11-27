package builtins

import (
	"compress/zlib"
	"errors"
	"io"
	"strconv"

	"github.com/alesiong/codec/codecs"
)

type zlibCodec struct {
}

func (b zlibCodec) RunCodec(input io.Reader, globalMode codecs.CodecMode, options map[string]string, output io.Writer) (err error) {
	level := -1
	if options["L"] != "" {
		level, err = strconv.Atoi(options["L"])
		if err != nil {
			return
		}
	}

	switch globalMode {
	case codecs.CodecModeEncoding:
		var zlibWriter *zlib.Writer
		zlibWriter, err = zlib.NewWriterLevel(output, level)
		if err != nil {
			return
		}
		_, err = io.Copy(zlibWriter, input)
		if err != nil {
			return
		}
		err = zlibWriter.Close()
		if err != nil {
			return
		}

		return
	case codecs.CodecModeDecoding:
		var zlibReader io.ReadCloser
		zlibReader, err = zlib.NewReader(input)
		if err != nil {
			return
		}
		_, err = io.Copy(output, zlibReader)
		if err != nil {
			return
		}

		err = zlibReader.Close()
		if err != nil {
			return
		}

		return
	default:
		return errors.New("invalid codec mode")
	}
}
