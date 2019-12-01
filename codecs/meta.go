package codecs

import (
	"io"
	"sort"
)

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

func (c *CodecMetaInfo) Register(name string, codec Codec) {
	c.codecsMap[name] = codec
}

func (c *CodecMetaInfo) Lookup(name string) Codec {
	return c.codecsMap[name]
}

func (c *CodecMetaInfo) AllCodecs() (chan string, chan Codec) {
	codecChan := make(chan Codec)
	nameChan := make(chan string)

	keys := make([]string, 0, len(c.codecsMap))
	for name := range c.codecsMap {
		keys = append(keys, name)
	}
	sort.Strings(keys)

	go func() {
		for _, name := range keys {
			nameChan <- name
			codecChan <- c.codecsMap[name]
		}
		close(nameChan)
		close(codecChan)
	}()

	return nameChan, codecChan
}
