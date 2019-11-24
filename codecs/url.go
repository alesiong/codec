package codecs

import (
	"bytes"
	"errors"
	"io"
	"net/url"
)

type urlCodecs struct {
}

func (h urlCodecs) Execute(input io.Reader, globalMode CodecMode, options map[string]string, output io.WriteCloser) (err error) {
	escape := url.QueryEscape
	unescape := url.QueryUnescape

	if options["p"] != "" {
		escape = url.PathEscape
		unescape = url.PathUnescape
	}

	switch globalMode {
	case CodecModeEncoding:
		err = ReadToWriter(input, &escapeWriter{
			escape: escape,
			writer: output,
		}, output)
	case CodecModeDecoding:
		err = ReadToWriter(input, &unescapeWriter{
			unescape: unescape,
			writer:   output,
		}, output)
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
	return e.writer.Write([]byte(e.escape(string(p))))
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
	return u.writer.Write([]byte(result))
}
