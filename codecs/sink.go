package codecs

import (
	"errors"
	"io"
	"os"
)

type teeCodecs struct {
}

func (t teeCodecs) Execute(input io.Reader, globalMode CodecMode, options map[string]string, output io.WriteCloser) (err error) {
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
	if len(writers) == 0 {
		return output.Close()
	}
	writer := io.MultiWriter(writers...)
	return ReadToWriter(input, writer, output)
}

type sinkCodecs struct {
	teeCodecs
}

func (s sinkCodecs) Execute(input io.Reader, globalMode CodecMode, options map[string]string, output io.WriteCloser) (err error) {
	return s.teeCodecs.Execute(input, globalMode, map[string]string{"c": "*"}, output)
}

type redirectCodecs struct {
	teeCodecs
}

func (r redirectCodecs) Execute(input io.Reader, globalMode CodecMode, options map[string]string, output io.WriteCloser) (err error) {
	outputFile := options["O"]
	if outputFile == "" {
		return errors.New("redirect: missing required option output file (-O)")
	}
	return r.teeCodecs.Execute(input, globalMode, map[string]string{
		"c": "*",
		"O": outputFile,
	}, output)
}
