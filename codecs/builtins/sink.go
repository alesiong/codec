package builtins

import (
	"errors"
	"io"
	"os"

	"github.com/alesiong/codec/codecs"
)

func init() {
	codecs.Register("tee", teeCodecs{})
	codecs.Register("sink", sinkCodecs{})
	codecs.Register("redirect", redirectCodecs{})
}

type teeCodecs struct {
}

func (t teeCodecs) RunCodec(input io.Reader, globalMode codecs.CodecMode, options map[string]string, output io.Writer) (err error) {
	writers := make([]io.Writer, 0, 2)
	if options["c"] == "" {
		writers = append(writers, output)
	}
	outputFile := options["O"]
	if outputFile != "" {
		var file *os.File
		file, err = os.OpenFile(outputFile, os.O_CREATE|os.O_WRONLY, 0755)
		if err != nil {
			return
		}
		defer file.Close()
		writers = append(writers, file)
	}
	writer := io.MultiWriter(writers...)
	_, err = io.Copy(writer, input)
	return
}

type sinkCodecs struct {
	teeCodecs
}

func (s sinkCodecs) RunCodec(input io.Reader, globalMode codecs.CodecMode, options map[string]string, output io.Writer) (err error) {
	return s.teeCodecs.RunCodec(input, globalMode, map[string]string{"c": "*"}, output)
}

type redirectCodecs struct {
	teeCodecs
}

func (r redirectCodecs) RunCodec(input io.Reader, globalMode codecs.CodecMode, options map[string]string, output io.Writer) (err error) {
	outputFile := options["O"]
	if outputFile == "" {
		return errors.New("redirect: missing required option output file (-O)")
	}
	return r.teeCodecs.RunCodec(input, globalMode, map[string]string{
		"c": "*",
		"O": outputFile,
	}, output)
}
