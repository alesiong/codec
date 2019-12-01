package builtins

import (
	"bytes"
	"errors"
	"io"
	"strconv"

	"github.com/alesiong/codec/codecs"
)

type escapeCodecs struct {
}

func init() {
	codecs.Register("escape", escapeCodecs{})
}

func (h escapeCodecs) RunCodec(input io.Reader, globalMode codecs.CodecMode, options map[string]string, output io.Writer) (err error) {
	escape := func(str string) string {
		result := strconv.Quote(str)
		return result[1 : len(result)-1]
	}
	unescape := func(s string) (string, error) {
		s = "\"" + s + "\""
		return strconv.Unquote(s)
	}

	switch globalMode {
	case codecs.CodecModeEncoding:
		_, err = io.Copy(&escapeWriter{
			escape: escape,
			writer: output,
		}, input)
	case codecs.CodecModeDecoding:
		_, err = io.Copy(output, &unquoteReader{
			unquote: unescape,
			reader:  input,
		})
	default:
		return errors.New("invalid codec mode")
	}

	return
}

type unquoteReader struct {
	unquote      func(string) (string, error)
	reader       io.Reader
	writeBuffer  []byte
	remainBuffer []byte
	buffer       [1024]byte
}

func (u *unquoteReader) Read(p []byte) (n int, err error) {
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

	lastSlashIndex := bytes.LastIndexByte(data, '\\')
	length := len(data)
	if lastSlashIndex >= 0 {
		_, _, _, e := strconv.UnquoteChar(string(data[lastSlashIndex:]), '"')
		if e == nil {
			u.writeBuffer = nil
		} else {
			if readErr == io.EOF {
				return readN, e
			}
			u.writeBuffer = make([]byte, length-lastSlashIndex)
			copy(u.writeBuffer, data[lastSlashIndex:])
			data = data[:lastSlashIndex]
		}
	} else {
		u.writeBuffer = nil
	}

	if len(data) == 0 {
		if u.writeBuffer != nil {
			readErr = io.ErrUnexpectedEOF
		}
		return 0, readErr
	}
	result, err := u.unquote(string(data))
	if err != nil {
		return 0, err
	}
	n = copy(p, result)
	u.remainBuffer = []byte(result)[n:]
	return
}
