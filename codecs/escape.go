package codecs

import (
	"bytes"
	"errors"
	"io"
	"strconv"
)

type escapeCodecs struct {
}

func (h escapeCodecs) RunCodec(input io.Reader, globalMode CodecMode, options map[string]string, output io.Writer) (err error) {
	escape := func(str string) string {
		result := strconv.Quote(str)
		return result[1 : len(result)-1]
	}
	unescape := func(s string) (string, error) {
		s = "\"" + s + "\""
		return strconv.Unquote(s)
	}

	switch globalMode {
	case CodecModeEncoding:
		err = ReadToWriter(input, &escapeWriter{
			escape: escape,
			writer: output,
		})
	case CodecModeDecoding:
		err = ReadToWriter(input, &unquoteWriter{
			unquote: unescape,
			writer:  output,
		})
	default:
		return errors.New("invalid codec mode")
	}

	return
}

type unquoteWriter struct {
	unquote func(string) (string, error)
	writer  io.Writer
	buffer  []byte
}

func (u *unquoteWriter) Write(p []byte) (n int, err error) {
	data := make([]byte, len(u.buffer), len(p)+len(u.buffer))
	copy(data, u.buffer)
	data = append(data, p...)

	lastSlashIndex := bytes.LastIndexByte(data, '\\')
	length := len(data)
	if lastSlashIndex >= 0 {
		_, _, _, e := strconv.UnquoteChar(string(data[lastSlashIndex:]), '"')
		if e == nil {
			u.buffer = nil
		} else {
			u.buffer = make([]byte, length-lastSlashIndex)
			copy(u.buffer, data[lastSlashIndex:])
			data = data[:lastSlashIndex]
		}
	} else {
		u.buffer = nil
	}

	if len(data) == 0 {
		return len(p), nil
	}
	result, err := u.unquote(string(data))
	if err != nil {
		return 0, err
	}
	return u.writer.Write([]byte(result))
}
