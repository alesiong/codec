package main

import (
	"errors"
	"io"

	"github.com/alesiong/codec/codecs"
)

type trCodecs struct {
}

func init() {
	codecs.Register("tr", trCodecs{})
}

func (t trCodecs) RunCodec(input io.Reader, globalMode codecs.CodecMode, options map[string]string, output io.Writer) (err error) {
	from := options["F"]
	if from == "" {
		return errors.New("tr: missing required option from string (-F)")
	}
	to := options["T"]
	if to == "" {
		return errors.New("tr: missing required option to string (-T)")
	}
	_, err = io.Copy(&trWriter{from, to, output}, input)
	return
}

type trWriter struct {
	from   string
	to     string
	writer io.Writer
}

func (t *trWriter) Write(p []byte) (n int, err error) {
	replacement := make(map[byte]byte, len(t.from))
	result := make([]byte, 0, len(p))
	for i := 0; i < len(t.from); i++ {
		if i < len(t.to) {
			replacement[t.from[i]] = t.to[i]
		} else {
			replacement[t.from[i]] = 0
		}
	}
	for _, c := range p {
		if r, ok := replacement[c]; ok {
			if r > 0 {
				result = append(result, r)
			}
		} else {
			result = append(result, c)
		}
	}
	_, err = t.writer.Write(result)
	return len(p), err
}

var CodecPlugin codecs.Codec = trCodecs{}

func main() {
}
