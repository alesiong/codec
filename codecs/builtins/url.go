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
		_, err = io.Copy(&escapeWriter{
			escape: escape,
			writer: output,
		}, input)
	case codecs.CodecModeDecoding:
		_, err = io.Copy(output, &unescapeReader{
			unescape: unescape,
			reader:   input,
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

type unescapeReader struct {
	unescape     func(string) (string, error)
	reader       io.Reader
	writeBuffer  []byte
	remainBuffer []byte
	buffer       [1024]byte
}

func (u *unescapeReader) Read(p []byte) (n int, err error) {
	if len(u.remainBuffer) > 0 {
		n = copy(p, u.remainBuffer)
		u.remainBuffer = u.remainBuffer[n:]
		return
	}

	readN, readErr := u.reader.Read(u.buffer[:])
	if readErr != nil && readErr != io.EOF {
		return readN, readErr
	}
	data := make([]byte, len(u.writeBuffer), readN+len(u.writeBuffer))
	copy(data, u.writeBuffer)
	data = append(data, u.buffer[:readN]...)

	lastPercentIndex := bytes.LastIndexByte(data, '%')
	length := len(data)
	if lastPercentIndex >= 0 && lastPercentIndex >= length-2 {
		u.writeBuffer = make([]byte, length-lastPercentIndex)
		copy(u.writeBuffer, data[lastPercentIndex:])
		data = data[:lastPercentIndex]
	} else {
		u.writeBuffer = nil
	}

	if len(data) == 0 {
		if u.writeBuffer != nil {
			readErr = io.ErrUnexpectedEOF
		}
		return 0, readErr
	}

	result, err := u.unescape(string(data))
	if err != nil {
		return 0, err
	}
	n = copy(p, result)
	u.remainBuffer = []byte(result)[n:]
	return
}
