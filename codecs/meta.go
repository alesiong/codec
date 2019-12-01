package codecs

import "io"

type MetaCodec interface {
	RunMetaCodec(input io.Reader, globalMode CodecMode, options map[string]string, metaInfo *CodecMetaInfo, output io.Writer) (err error)
}

type metaCodecWrap struct {
	metaCodec MetaCodec
}

func (m metaCodecWrap) RunCodec(input io.Reader, globalMode CodecMode, options map[string]string, output io.Writer) (err error) {
	return m.metaCodec.RunMetaCodec(input, globalMode, options, &codecMetaInfo, output)
}

func RegisterMeta(name string, metaCodec MetaCodec) {
	Register(name, metaCodecWrap{metaCodec})
}
