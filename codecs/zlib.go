package codecs

import (
	"compress/zlib"
	"errors"
	"io"
	"strconv"
)

type zlibCodec struct {
}

func (b zlibCodec) RunCodec(input io.Reader, globalMode CodecMode, options map[string]string, output io.Writer) (err error) {
	level := -1
	if options["L"] != "" {
		level, err = strconv.Atoi(options["L"])
		if err != nil {
			return
		}
	}

	switch globalMode {
	case CodecModeEncoding:
		var zlibWriter *zlib.Writer
		zlibWriter, err = zlib.NewWriterLevel(output, level)
		if err != nil {
			return
		}
		err = ReadToWriter(input, zlibWriter)
		if err != nil {
			return
		}
		err = zlibWriter.Close()
		if err != nil {
			return
		}

		return
	case CodecModeDecoding:
		var zlibReader io.ReadCloser
		zlibReader, err = zlib.NewReader(input)
		if err != nil {
			return
		}
		err = ReadToWriter(zlibReader, output)
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
