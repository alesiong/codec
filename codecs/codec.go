package codecs

import (
	"io"
)

type CodecMode int

const (
	CodecModeEncoding CodecMode = iota
	CodecModeDecoding
)

type Codec interface {
	RunCodec(input io.Reader, globalMode CodecMode, options map[string]string, output io.Writer) (err error)
}
