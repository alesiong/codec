package codecs

import (
	"io"
	"sync"
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

func (c *CodecMetaInfo) Register(name string, codec Codec) {
	c.codecsMap[name] = codec
}

func (c *CodecMetaInfo) Lookup(name string) Codec {
	return c.codecsMap[name]
}

func (c *CodecMetaInfo) AllCodecs() (chan string, chan Codec) {
	codecChan := make(chan Codec)
	nameChan := make(chan string)
	var wg sync.WaitGroup

	for name, codec := range c.codecsMap {
		go func(name string, codec Codec) {
			nameChan <- name
			codecChan <- codec
			wg.Done()
		}(name, codec)
		wg.Add(1)
	}
	go func() {
		wg.Wait()
		close(nameChan)
		close(codecChan)
	}()
	return nameChan, codecChan
}

func Register(name string, codec Codec) {
	codecMetaInfo.Register(name, codec)
}

func Lookup(name string) Codec {
	return codecMetaInfo.Lookup(name)
}
