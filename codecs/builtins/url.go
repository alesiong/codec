package builtins

import (
	"bytes"
	"errors"
	"io"
	"net/url"

	"github.com/alesiong/codec/codecs"
)

type urlCodecs struct {
}

func (h urlCodecs) RunCodec(input io.Reader, globalMode codecs.CodecMode, options map[string]string, output io.Writer) (err error) {
	escape := url.QueryEscape
	unescape := url.QueryUnescape

	if options["p"] != "" {
		escape = url.PathEscape
		unescape = url.PathUnescape
	}

	switch globalMode {
	case codecs.CodecModeEncoding:
		err = codecs.ReadToWriter(input, &escapeWriter{
			escape: escape,
			writer: output,
		})
	case codecs.CodecModeDecoding:
		err = codecs.ReadToWriter(input, &unescapeWriter{
			unescape: unescape,
			writer:   output,
		})
	default:
		return errors.New("invalid codec mode")
	}

	return
}

type escapeWriter struct {
	escape func(string) string
	writer io.Writer
}

func (e *escapeWriter) Write(p []byte) (n int, err error) {
	n = len(p)
	_, err = e.writer.Write([]byte(e.escape(string(p))))
	return
}

type unescapeWriter struct {
	unescape func(string) (string, error)
	writer   io.Writer
	buffer   []byte
}

func (u *unescapeWriter) Write(p []byte) (n int, err error) {
	data := make([]byte, len(u.buffer), len(p)+len(u.buffer))
	copy(data, u.buffer)
	data = append(data, p...)

	lastPercentIndex := bytes.LastIndexByte(data, '%')
	length := len(data)
	if lastPercentIndex >= 0 && lastPercentIndex >= length-2 {
		u.buffer = make([]byte, length-lastPercentIndex)
		copy(u.buffer, data[lastPercentIndex:])
		data = data[:lastPercentIndex]
	} else {
		u.buffer = nil
	}

	if len(data) == 0 {
		return len(p), nil
	}

	result, err := u.unescape(string(data))
	if err != nil {
		return 0, err
	}
	n = len(p)
	_, err = u.writer.Write([]byte(result))
	return
}
