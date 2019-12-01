package codecs

import (
	"io"
)

type CodecMode int

const (
	CodecModeEncoding CodecMode = iota
	CodecModeDecoding
)

var (
	codecMetaInfo = CodecMetaInfo{codecsMap: make(map[string]Codec)}
)

type Codec interface {
	RunCodec(input io.Reader, globalMode CodecMode, options map[string]string, output io.Writer) (err error)
}

type CodecUsage interface {
	Usage() string
}

type CodecMetaInfo struct {
	codecsMap map[string]Codec
}

func Register(name string, codec Codec) {
	codecMetaInfo.Register(name, codec)
}

func Lookup(name string) Codec {
	return codecMetaInfo.Lookup(name)
}
